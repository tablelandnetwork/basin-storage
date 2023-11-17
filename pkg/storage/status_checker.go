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
	// StatusClient is a w3s.Client instance used to interact with W3S.
	StatusClient w3s.Client
	// DBClient is a Crdb instance used to interact with CockroachDB.
	DBClient Crdb
	// contractClient is a BasinStorage contract interface
	contractClient ethereum.BasinStorage
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

// addCID prepares Tx to add a CID to the contract.
func (sc *StatusChecker) addCID(
	ctx context.Context,
	pub string,
	cid string,
	timestamp int64,
) error {
	// prepare tx opts with gas related params
	txOpts, err := sc.contractClient.EstimateGas(ctx, pub, cid, timestamp)
	if err != nil {
		return fmt.Errorf("failed to estimate gas for adding cid: %v", err)
	}

	// set nonce
	nonce, err := sc.contractClient.GetPendingNonce(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(nonce))

	fmt.Println("Adding cid: ", pub, cid, timestamp, nonce)
	if err = sc.contractClient.AddCID(ctx, pub, cid, timestamp, txOpts); err != nil {
		return fmt.Errorf("failed to add cid to contract: %v", err)
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

// checkActiveDeals returns true if there are any
// active deals for the job.
func (sc *StatusChecker) checkActiveDeals(
	status *w3s.Status,
	job UnfinishedJob,
) bool {
	// when there are no deals returned by W3S
	if len(status.Deals) == 0 {
		fmt.Printf(
			"no deals found for job: %s, %x \n",
			job.Pub, job.Cid)
		return false
	}

	// when deals exist, check if they are active
	deals := []w3s.Deal{}
	for _, d := range takeActiveDeals(status.Deals) {
		// filter out deals that are not active yet
		fmt.Printf("deal status: %s \n", d.Status)
		deals = append(deals, d)
	}
	if len(deals) == 0 {
		fmt.Printf(
			"deals exist, but are not activated: %s, %x \n",
			job.Pub, job.Cid)
		return false
	}

	return true
}

func (sc *StatusChecker) processJob(
	ctx context.Context,
	job UnfinishedJob,
) error {
	fmt.Printf("checking status for job: %s, %x\n", job.Pub, job.Cid)
	pub := fmt.Sprintf("%s.%s", job.Pub.Namespace, job.Pub.Relation)

	// Check Job status
	status, err := sc.getStatus(ctx, job.Cid)
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}

	// Find active activeDeals for the jobs
	activeDeals := sc.checkActiveDeals(status, job)
	if !activeDeals {
		fmt.Println("skipping indexing cid")
		return nil
	}

	cid, err := cid.Cast(job.Cid)
	if err != nil {
		return fmt.Errorf("failed to cast cid from bytes: %v", err)
	}

	var ts int64
	if job.Timestamp == nil {
		ts = int64(0)
	} else {
		ts = *job.Timestamp
	}

	if err := sc.addCID(ctx, pub, cid.String(), ts); err != nil {
		return fmt.Errorf("failed to add cid: %v", err)
	}

	return sc.updateJobStatus(ctx, job, status)
}

// ProcessJobs checks the status of all unfinished jobs.
// If a job has active deals, it adds the "CID" to the BasinStorage contract.
// If a job has no active deals, it does nothing.
// If a job has already been activated, it does nothing.
// Finally, it updates the job status in the DB.
func (sc *StatusChecker) ProcessJobs(ctx context.Context) error {
	unfinishedJobs, err := sc.DBClient.UnfinishedJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get unfinished jobs: %v", err)
	}

	for _, job := range unfinishedJobs {
		if err := sc.processJob(ctx, job); err != nil {
			return fmt.Errorf("failed to process job: %v", err)
		}
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
