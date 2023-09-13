package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
)

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
				Activated: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Pub:       "myfile2",
				Cid:       getCIDFromBytes([]byte("data for myfile2")).Bytes(),
				Activated: time.Time{}, // not marked as active but deals are active on chain
			},
			{
				Pub:       "myfile3",
				Cid:       getCIDFromBytes([]byte("data for myfile3")).Bytes(),
				Activated: time.Time{}, // not marked as active and deals are in queue
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

	assert.Equal(t, 2, len(bsc.deals))
	assert.Equal(t, uint64(3), bsc.deals[0].Id)
	assert.Equal(t, uint64(4), bsc.deals[1].Id)

	var ts time.Time
	for _, j := range db.jobs {
		if j.Pub == "myfile" {
			ts = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
			assert.Equal(t, ts, j.Activated)
		}

		if j.Pub == "myfile2" {
			ts = time.Date(2021, time.January, 5, 3, 0, 0, 0, time.UTC)
			assert.Equal(t, ts, j.Activated)
		}

		if j.Pub == "myfile3" {
			ts = time.Time{}
			assert.Equal(t, ts, j.Activated)
		}
	}
}
