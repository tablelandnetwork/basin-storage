package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"io/fs"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"

	w3s "github.com/web3-storage/go-w3s-client"
	w3http "github.com/web3-storage/go-w3s-client/http"
)

// Mock interface for w3s.Client.
type mockW3sClient struct {
	Files []fs.File
}

func mockData() []byte {
	return []byte("hello world")
}

func getCIDFromBytes(mockData []byte) cid.Cid {
	hashedData := sha256.Sum256(mockData)
	multihash, _ := mh.Encode(hashedData[:], mh.SHA2_256)
	cidV1 := cid.NewCidV1(cid.Raw, multihash)
	return cidV1
}

func (m *mockW3sClient) Put(_ context.Context, file fs.File, _ ...w3s.PutOption) (cid.Cid, error) {
	m.Files = append(m.Files, file)
	return getCIDFromBytes(mockData()), nil
}

func (m *mockW3sClient) Get(_ context.Context, _ cid.Cid) (*w3http.Web3Response, error) {
	return nil, nil
}

func (m *mockW3sClient) PutCar(_ context.Context, _ io.Reader) (cid.Cid, error) {
	return cid.Cid{}, nil
}

// mock deals that are active on chain and in db.
var activeDealsJob1 = []w3s.Deal{
	{
		Activation:        time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		DealID:            1,
		DataModelSelector: "foo/bar",
		Status:            w3s.DealStatusActive,
	},
	{
		Activation:        time.Date(2021, time.January, 1, 3, 0, 0, 0, time.UTC),
		DealID:            2,
		DataModelSelector: "foo/bar",
		Status:            w3s.DealStatusActive,
	},
}

// mock deals that are active on chain but not in db.
var activeDealsJob2 = []w3s.Deal{
	// same deal id as job1 but different cid and selector
	// it may happen that two different files with different CIDs
	// gets added to the same deal.
	{
		Activation:        time.Date(2021, time.January, 5, 3, 0, 0, 0, time.UTC),
		DealID:            1,
		DataModelSelector: "bar/foo",
		Status:            w3s.DealStatusActive,
	},
	{
		Activation:        time.Date(2021, time.January, 5, 5, 0, 0, 0, time.UTC),
		DealID:            4,
		DataModelSelector: "bar/foo",
		Status:            w3s.DealStatusActive,
	},
}

// mock deals that are not active on chain.
var inactiveDealsJob3 = []w3s.Deal{
	{
		Activation:        time.Date(2021, time.January, 7, 4, 0, 0, 0, time.UTC),
		DealID:            5,
		DataModelSelector: "baz/foo",
		Status:            w3s.DealStatusQueued,
	},
	{
		Activation:        time.Date(2021, time.January, 7, 6, 0, 0, 0, time.UTC),
		DealID:            6,
		DataModelSelector: "baz/foo",
		Status:            w3s.DealStatusPublished,
	},
}

func (m *mockW3sClient) Status(_ context.Context, cid cid.Cid) (*w3s.Status, error) {
	var deals []w3s.Deal
	cid1 := getCIDFromBytes([]byte("data for myfile"))
	cid2 := getCIDFromBytes([]byte("data for myfile2"))

	// job 1 has active deals
	if bytes.Equal(cid.Bytes(), cid1.Bytes()) {
		deals = activeDealsJob1
	} else if bytes.Equal(cid.Bytes(), cid2.Bytes()) {
		// job2 has active deals but not marked in db
		deals = activeDealsJob2
	} else {
		deals = inactiveDealsJob3
	}

	status := &w3s.Status{
		Cid:     cid,
		DagSize: 0,
		Created: time.Now(),
		Pins:    nil,
		Deals:   deals,
	}
	return status, nil
}

func (m *mockW3sClient) List(_ context.Context, _ ...w3s.ListOption) (*w3s.UploadIterator, error) {
	return nil, nil
}

func (m *mockW3sClient) Pin(_ context.Context, _ cid.Cid, _ ...w3s.PinOption) (*w3s.PinResponse, error) {
	return nil, nil
}

// MockReadCloser is a mock type for crdb.DBClient.
type MockReadCloser struct {
	*bytes.Reader
}

// Close is a mock implementation of io.Closer.
func (mrc *MockReadCloser) Close() error {
	return nil
}

type mockCrdb struct {
	jobs []UnfinishedJob
}

func (m *mockCrdb) CreateJob(
	_ context.Context,
	cidStr string,
	fname string,
	timestamp *int64,
	cacheDuration int64,
	_ string,
	_ string,
) error {
	cid, _ := cid.Decode(cidStr)
	pub, err := extractPub(fname)
	if err != nil {
		return err
	}
	m.jobs = append(m.jobs, UnfinishedJob{
		Pub:       pub,
		Cid:       cid.Bytes(),
		Activated: time.Time{},
		Timestamp: timestamp,
		ExpiresAt: time.Unix(*timestamp+cacheDuration, 0),
	})
	return nil
}

func (m *mockCrdb) UnfinishedJobs(_ context.Context) ([]UnfinishedJob, error) {
	var t time.Time
	ufj := []UnfinishedJob{}
	for _, job := range m.jobs {
		if job.Activated == t {
			ufj = append(ufj, job)
		}
	}
	return ufj, nil
}

func (m *mockCrdb) UpdateJobStatus(_ context.Context, cid []byte, activation time.Time) error {
	for i, job := range m.jobs {
		if bytes.Equal(job.Cid, cid) {
			m.jobs[i].Activated = activation
			break
		}
	}
	return nil
}

// MockBasinStorage is the mock type for BasinStorage Contract.
type MockBasinStorage struct {
	cids []string
}

// EstimateGas is a mock implementation of BasinStorage.EstimateGas.
func (c *MockBasinStorage) EstimateGas(
	_ context.Context,
	_ string,
	_ string,
	_ int64,
) (*bind.TransactOpts, error) {
	return &bind.TransactOpts{}, nil
}

// GetPendingNonce is a mock implementation of BasinStorage.GetPendingNonce.
func (c *MockBasinStorage) GetPendingNonce(
	_ context.Context,
) (uint64, error) {
	return 0, nil
}

// AddCID is a mock implementation of BasinStorage.AddCID.
func (c *MockBasinStorage) AddCID(
	_ context.Context,
	_ string,
	cids string,
	_ int64,
	_ *bind.TransactOpts,
) error {
	time.Sleep(1 * time.Second) // fake delay
	c.cids = append(c.cids, cids)
	return nil
}
