package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatusChecker(t *testing.T) {
	ctx := context.Background()
	bsc := &MockBasinStorage{
		cids: []string{},
	}
	db := &mockCrdb{
		jobs: []UnfinishedJob{
			{
				Pub: Pub{Namespace: "testns", Relation: "testrel"},
				Cid: getCIDFromBytes([]byte("data for myfile")).Bytes(),
				// marked as active and deals are active on chain
				// CID should be skipped
				Activated: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Pub: Pub{Namespace: "testns2", Relation: "testrel2"},
				Cid: getCIDFromBytes([]byte("data for myfile2")).Bytes(),
				// not marked as active but deals are active on chain
				// CID should be added
				Activated: time.Time{},
			},
			{
				Pub: Pub{Namespace: "testns", Relation: "testrel3"},
				Cid: getCIDFromBytes([]byte("data for myfile3")).Bytes(),
				// not marked as active and deals are in queue
				// CID cannot be added
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

	expectedCidStr := getCIDFromBytes([]byte("data for myfile2")).String()

	assert.Equal(t, 1, len(bsc.cids))		
	assert.Equal(t, expectedCidStr, bsc.cids[0])
		

	var ts time.Time
	for _, j := range db.jobs {
		if j.Pub.Relation == "testrel" {
			ts = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
			assert.Equal(t, ts, j.Activated)
		}

		if j.Pub.Relation == "restrel2" {
			ts = time.Date(2021, time.January, 5, 3, 0, 0, 0, time.UTC)
			assert.Equal(t, ts, j.Activated)
		}

		if j.Pub.Relation == "testRel3" {
			ts = time.Time{}
			assert.Equal(t, ts, j.Activated)
		}
	}
}
