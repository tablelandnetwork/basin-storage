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
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
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

// mock deals that are active on chain and in db
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

// mock deals that are active on chain but not in db
// (status checking job has not run yet)
var activeDealsJob2 = []w3s.Deal{
	{
		Activation:        time.Date(2021, time.January, 5, 3, 0, 0, 0, time.UTC),
		DealID:            3,
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

// mock deals that are not active on chain
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

	// job 1 // has active deals
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

type MockReadCloser struct {
	*bytes.Reader
}

func (mrc *MockReadCloser) Close() error {
	return nil
}

type mockCrdb struct {
	jobs []UnfinihedJob
}

func (m *mockCrdb) CreateJob(_ context.Context, cidStr string, pub string) error {
	cid, _ := cid.Decode(cidStr)
	m.jobs = append(m.jobs, UnfinihedJob{
		Pub:       pub,
		Cid:       cid.Bytes(),
		Activated: time.Time{},
	})
	return nil
}

func (m *mockCrdb) UnfinishedJobs(_ context.Context) ([]UnfinihedJob, error) {
	var t time.Time
	ufj := []UnfinihedJob{}
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

// Mock interface for BasinStorage Contract.
type MockBasinStorage struct {
	deals []ethereum.BasinStorageDealInfo
}

func (c *MockBasinStorage) EstimateGas(
	_ context.Context,
	_ *bind.TransactOpts,
	_ string,
	_ []ethereum.BasinStorageDealInfo,
) (*bind.TransactOpts, error) {
	return &bind.TransactOpts{}, nil
}

func (c *MockBasinStorage) GetRecentDeals(
	_ context.Context, _ string,
) (map[ethereum.BasinStorageDealInfo]struct{}, error) {
	return nil, nil
}

func (c *MockBasinStorage) AddDeals(
	_ context.Context,
	_ string,
	deals []ethereum.BasinStorageDealInfo,
) error {
	time.Sleep(1 * time.Second) // fake delay
	c.deals = append(c.deals, deals...)
	return nil
}
