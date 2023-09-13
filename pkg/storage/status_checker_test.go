package storage

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
)

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
	c.deals = append(c.deals, deals...)
	return nil
}

func TestStatusChecker(t *testing.T) {
	ctx := context.Background()

	bsc := &MockBasinStorage{
		deals: []ethereum.BasinStorageDealInfo{},
	}
	db := &mockCrdb{
		jobs: []UnfinihedJob{
			{
				Pub:       "myfile",
				Cid:       getCIDFromBytes([]byte("data for myfile")).Bytes(),
				Activated: time.Now().Add(-time.Hour * 2), // activated 2 hours ago
			},
			{
				Pub:       "myfile2",
				Cid:       getCIDFromBytes([]byte("data for myfile2")).Bytes(),
				Activated: time.Time{},
			},
		},
	}

	sc := StatusChecker{
		StatusClient:   &mockW3sClient{},
		DBClient:       db,
		contractClient: bsc,
	}
	err := sc.ProcessJobs(ctx)
	assert.NoError(t, err)

	for _, j := range db.jobs {
		// both (mock) jobs should be activated
		// after the call to ProcessJobs
		assert.NotEqual(t, time.Time{}, j.Activated)
	}
}
