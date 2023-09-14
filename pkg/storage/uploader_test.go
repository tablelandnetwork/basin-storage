package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUploader(t *testing.T) {
	ctx := context.Background()
	mockGCS := new(mocks.GCS)

	// Mocking the returned values for the ParseEventData method
	fname := "feeds_yyyy-mm-dd-000-111-1-2-00000000-basin_storage.employees.employees-2.parquet"
	mockGCS.On("ParseEvent").Return("mybucket", fname, nil)

	// Mocking the returned reader for the GetObjectReader method
	mockReadCloser := &MockReadCloser{Reader: bytes.NewReader(mockData())}
	mockGCS.On("GetObjectReader", ctx, "mybucket", fname).Return(mockReadCloser, nil)

	uploader := FileUploader{
		StorageClient: mockGCS,
		DealClient: &mockW3sClient{
			Files: []fs.File{},
		},
		DBClient: &mockCrdb{
			jobs: []UnfinihedJob{},
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

	// Assert that the CID was added to the database with the correct pub name
	db := uploader.DBClient.(*mockCrdb)
	assert.Equal(t, 1, len(db.jobs))
	fmt.Println(db.jobs)
	expectedPub := Pub{Namespace: "employees", Relation: "employees"}
	assert.Equal(t, expectedPub, db.jobs[0].Pub)

	cid, err := cid.Parse(db.jobs[0].Cid)
	assert.NoError(t, err)
	assert.Equal(t, getCIDFromBytes(mockData()).String(), cid.String())
}
