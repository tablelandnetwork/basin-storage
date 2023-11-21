package storage

import (
	"bytes"
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUploader(t *testing.T) {
	ctx := context.Background()
	mockGCS := new(mocks.GCS)

	// Mocking the returned values for the ParseEventData method
	fname := "foo.bar.baz/relname/exportabcd1234-2.0.parquet"
	mockGCS.On("ParseEvent").Return("mybucket", fname, nil)

	// Mocking the returned reader for the GetObjectReader method
	mockReadCloser := &MockReadCloser{Reader: bytes.NewReader(mockData())}
	mockGCS.On("GetObjectReader", ctx, "mybucket", fname).Return(mockReadCloser, nil)
	mockGCS.On("GetObjectMetadata", ctx, "mybucket", fname).Return(
		map[string]string{"timestamp": "1700248832", "cache_duration": "100"}, nil,
	)

	uploader := FileUploader{
		StorageClient: mockGCS,
		DealClient: &mockW3sClient{
			Files: []fs.File{},
		},
		DBClient: &mockCrdb{
			jobs: []UnfinishedJob{},
		},
	}

	err := uploader.Upload(ctx)

	// Assert that mockGCS.GetObjectReader was called with the correct arguments
	mockGCS.AssertExpectations(t)
	assert.NoError(t, err)

	files := uploader.DealClient.(*mockW3sClient).Files
	assert.Equal(t, 1, len(files))

	fStat, err := files[0].Stat()
	assert.NoError(t, err)
	assert.Equal(t, fname, fStat.Name())
	assert.Equal(t, int64(11), fStat.Size())

	// Read the file contents that was send to (mocked) w3s
	buf := make([]byte, 11)
	_, err = files[0].Read(buf)
	assert.NoError(t, err)

	err = files[0].Close()
	assert.NoError(t, err)

	assert.Equal(t, mockData(), buf)

	jobs, err := uploader.DBClient.UnfinishedJobs(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))

	expectedPub := Pub{Namespace: "foo.bar.baz", Relation: "relname"}
	assert.Equal(t, expectedPub, jobs[0].Pub)

	cid, err := cid.Parse(jobs[0].Cid)
	assert.NoError(t, err)
	assert.Equal(t, getCIDFromBytes(mockData()).String(), cid.String())

	assert.Equal(t, int64(1700248832), *jobs[0].Timestamp)
	assert.Equal(t, time.Unix(1700248832+100, 0), jobs[0].ExpiresAt)
}
