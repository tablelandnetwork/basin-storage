package storage

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
	"github.com/textileio/go-tableland/pkg/wallet"

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
	StatusClient   w3s.Client            // DealClient is a w3s.Client instance used to interact with W3S.
	DBClient       Crdb                  // DBClient is a Crdb instance used to interact with CockroachDB.
	contractClient ethereum.BasinStorage // *ethereum.Client // TODO: change to an interface for testing
}

// NewStatusChecker creates a new StatusChecker.
func NewStatusChecker(ctx context.Context, cfg *StatusCheckerConfig) (*StatusChecker, error) {
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

func (sc *StatusChecker) findEarliestDeal(deals []w3s.Deal) w3s.Deal {
	earliestDeal := deals[0]

	for _, d := range deals {
		if d.Activation.Before(earliestDeal.Activation) {
			earliestDeal = d
		}
	}

	return earliestDeal
}

func (sc *StatusChecker) processJob(ctx context.Context, job UnfinihedJob) error {
	fmt.Printf("checking status for job: %s, %x\n", job.Pub, job.Cid)
	status, err := sc.getStatus(ctx, job.Cid)
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}

	if len(status.Deals) == 0 {
		fmt.Printf("no deals found for job, skipping: %s, %x \n", job.Pub, job.Cid)
		return nil
	}

	deals := []ethereum.BasinStorageDealInfo{}
	for _, d := range status.Deals {
		// filter out deals that are not active yet
		if d.Status == w3s.DealStatusActive {
			deals = append(deals, ethereum.BasinStorageDealInfo{
				Id:           d.DealID,
				SelectorPath: d.DataModelSelector,
			})
		}
	}

	if len(deals) == 0 {
		fmt.Printf(
			"exitsing deals for the job are not activated, skipping: %s, %x \n",
			job.Pub, job.Cid)
		return nil
	}

	// Add deals to the BasinStorage contract
	if err = sc.contractClient.AddDeals(ctx, job.Pub, deals); err != nil {
		return fmt.Errorf("failed to add deals to contract: %v", err)
	}

	// Update job status in DB
	firstDealTS := sc.findEarliestDeal(status.Deals).Activation
	if err = sc.DBClient.UpdateJobStatus(ctx, job.Cid, firstDealTS); err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	fmt.Printf("finished updating status for job: %s, %x \n", job.Pub, job.Cid)
	return nil
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

	// asnychronously process unfinished jobs
	errs, ctx := errgroup.WithContext(ctx)
	for _, job := range unfinihedJobs {
		ctx := ctx
		job := job
		errs.Go(func() error {
			return sc.processJob(ctx, job)
		})
	}

	if err := errs.Wait(); err != nil {
		return fmt.Errorf("one or more jobs failed: %v", err)
	}

	return nil
}
