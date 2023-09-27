package storage

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
	"github.com/textileio/go-tableland/pkg/wallet"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	w3s "github.com/web3-storage/go-w3s-client"
)

// StatusCheckerConfig defines the configuration for a StatusChecker.
type StatusCheckerConfig struct {
	W3SToken         string
	CrdbConn         string
	PrivateKey       string
	BackendURL       string
	BasinStorageAddr string
	ChainID          string
}

// StatusChecker checks the status of a job and updates the status in the DB.
type StatusChecker struct {
	// DealClient is a w3s.Client instance used to interact with W3S.
	StatusClient w3s.Client
	// DBClient is a Crdb instance used to interact with CockroachDB.
	DBClient Crdb
	// contractClient is a BasinStorage contract interface
	contractClient ethereum.BasinStorage
	// simulated run flag (for testing by triggering the status checker)
	simulated bool
}

// NewStatusChecker creates a new StatusChecker.
func NewStatusChecker(ctx context.Context, cfg *StatusCheckerConfig, sim bool) (*StatusChecker, error) {
	wallet, err := wallet.NewWallet(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wallet: %v", err)
	}

	backend, err := ethclient.DialContext(ctx, cfg.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize backend: %v", err)
	}

	addr, err := common.NewMixedcaseAddressFromString(cfg.BasinStorageAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to read basin storage address: %v", err)
	}

	chainID, err := strconv.ParseUint(cfg.ChainID, int(10), 64)
	if err != nil {
		return nil, fmt.Errorf("failed to read chain ID: %v", err)
	}

	ethClient, err := ethereum.NewClient(
		backend,
		backend,
		chainID,
		addr.Address(),
		wallet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethereum client: %v", err)
	}

	// Initialize web3.storage client to upload file
	w3sOpts := []w3s.Option{
		w3s.WithToken(cfg.W3SToken),
		w3s.WithHTTPClient(
			&http.Client{
				Timeout: 0, // no timeout
			},
		),
	}
	w3sClient, err := w3s.NewClient(w3sOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize web3.storage client: %v", err)
	}

	// Initialize cockroachdb client to store metadata
	dbClient, err := NewDB(cfg.CrdbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db client: %v", err)
	}

	return &StatusChecker{
		StatusClient:   w3sClient,
		DBClient:       dbClient,
		contractClient: ethClient,
		simulated:      sim,
	}, nil
}

func (sc *StatusChecker) getStatus(ctx context.Context, CIDBytes []byte) (*w3s.Status, error) {
	jobCid, err := cid.Parse(CIDBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cid: %v", err)
	}
	status, err := sc.StatusClient.Status(ctx, jobCid)
	if err != nil {
		return nil, fmt.Errorf("failed to call w3s: %v", err)
	}
	return status, nil
}

// addDeals prepares Tx and adds deals to the contract.
func (sc *StatusChecker) addDeals(
	ctx context.Context,
	pub string,
	deals []ethereum.BasinStorageDealInfo,
	nonce uint64,
) error {
	// prepare tx opts with gas related params
	txOpts, err := sc.contractClient.EstimateGas(ctx, pub, deals)
	if err != nil {
		return fmt.Errorf("failed to estimate gas for adding deals: %v", err)
	}
	// set nonce
	txOpts.Nonce = big.NewInt(int64(nonce))
	fmt.Println("Adding deals with nonce: ", nonce)
	if err = sc.contractClient.AddDeals(ctx, pub, deals, txOpts); err != nil {
		return fmt.Errorf("failed to add deals to contract: %v", err)
	}

	return nil
}

// updateJobStatus Updates job status in DB.
func (sc *StatusChecker) updateJobStatus(
	ctx context.Context,
	job UnfinishedJob,
	status *w3s.Status,
) error {
	fdts := findEarliestDeal(status.Deals).Activation
	if err := sc.DBClient.UpdateJobStatus(ctx, job.Cid, fdts); err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}
	fmt.Printf("finished updating status for job: %s, %x \n", job.Pub, job.Cid)
	return nil
}

// getActiveDeals returns active deals for a job.
func (sc *StatusChecker) getActiveDeals(
	status *w3s.Status,
	job UnfinishedJob,
) []ethereum.BasinStorageDealInfo {
	deals := []ethereum.BasinStorageDealInfo{}
	// when there are no deals returned by W3S
	if len(status.Deals) == 0 {
		fmt.Printf(
			"no deals found for job: %s, %x \n",
			job.Pub, job.Cid)
		return deals
	}
	for _, d := range takeActiveDeals(status.Deals) {
		// filter out deals that are not active yet
		deals = append(deals, ethereum.BasinStorageDealInfo{
			Id:           d.DealID,
			SelectorPath: d.DataModelSelector,
		})
	}
	if len(deals) == 0 {
		fmt.Printf(
			"deals exist, but are not activated: %s, %x \n",
			job.Pub, job.Cid)
	}
	return deals
}

func (sc *StatusChecker) processJob(
	ctx context.Context,
	job UnfinishedJob,
	nonce uint64,
) error {
	fmt.Printf("checking status for job: %s, %x\n", job.Pub, job.Cid)
	pub := fmt.Sprintf("%s.%s", job.Pub.Namespace, job.Pub.Relation)

	// Check Job status
	status, err := sc.getStatus(ctx, job.Cid)
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}

	// Find active deals for the jobs
	deals := sc.getActiveDeals(status, job)
	if len(deals) == 0 {
		fmt.Println("skipping adding deals")
		return nil
	}

	if !sc.simulated {
		fmt.Printf(
			"not simulated: checking for duplicates: %s, %x \n",
			pub, job.Cid)
		deals, err = sc.removeDuplicateDeals(ctx, pub, deals)
		// Don't fail if we can't remove duplicates, log it
		// and continue to add deals
		if err != nil {
			fmt.Println("failed to remove duplicates: ", err)
		}
	}

	if len(deals) > 0 {
		if err := sc.addDeals(ctx, pub, deals, nonce); err != nil {
			return fmt.Errorf("failed to add deals: %v", err)
		}
	} else {
		fmt.Printf(
			"skipping: all deals are already indexed: %s, %x \n",
			pub, job.Cid)
	}

	return sc.updateJobStatus(ctx, job, status)
}

// ProcessJobs checks the status of all unfinished jobs.
// If a job has deals, it adds the deals to the BasinStorage contract.
// If a job has no deals, it does nothing.
// If a job has already been activated, it does nothing.
// Finally, it updates the job status in the DB.
func (sc *StatusChecker) ProcessJobs(ctx context.Context) error {
	unfinihedJobs, err := sc.DBClient.UnfinishedJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get unfinished jobs: %v", err)
	}

	// get current nonce
	nonce, err := sc.contractClient.GetPendingNonce(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	errs, ctx := errgroup.WithContext(ctx)
	for idx, job := range unfinihedJobs {
		idx := idx
		ctx := ctx
		job := job
		errs.Go(func() error {
			return sc.processJob(ctx, job, nonce+uint64(idx))
		})
	}

	if err := errs.Wait(); err != nil {
		return fmt.Errorf("one or more jobs failed: %v", err)
	}

	return nil
}

func findEarliestDeal(deals []w3s.Deal) w3s.Deal {
	earliestDeal := deals[0]
	for _, d := range deals {
		if d.Activation.Before(earliestDeal.Activation) {
			earliestDeal = d
		}
	}
	return earliestDeal
}

func takeActiveDeals(deals []w3s.Deal) []w3s.Deal {
	activeDeals := []w3s.Deal{}

	for _, d := range deals {
		if d.Status == w3s.DealStatusActive {
			activeDeals = append(activeDeals, d)
		}
	}

	return activeDeals
}

func (sc *StatusChecker) removeDuplicateDeals(
	ctx context.Context, pub string,
	deals []ethereum.BasinStorageDealInfo,
) ([]ethereum.BasinStorageDealInfo, error) {
	// filter out deals that are already in the contract
	// by looking at the recent deals.
	// due to db failure and retries, it may happen that deals are already added
	// to avoid duplicates we filter out deals that are already in the contract
	recentDeals, err := sc.contractClient.GetRecentDeals(ctx, pub)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent deals: %v", err)
	}

	deDupDeals := []ethereum.BasinStorageDealInfo{}
	for _, d := range deals {
		if _, ok := recentDeals[d.Id]; !ok {
			deDupDeals = append(deDupDeals, d)
		} else {
			fmt.Println(
				"deal already exists, skipping",
				d.Id, d.SelectorPath)
		}
	}

	return deDupDeals, nil
}
