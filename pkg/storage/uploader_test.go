package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/mocks"
	w3s "github.com/web3-storage/go-w3s-client"
	w3http "github.com/web3-storage/go-w3s-client/http"

	mh "github.com/multiformats/go-multihash"

	"github.com/stretchr/testify/assert"
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

func (m *mockW3sClient) Status(_ context.Context, cid cid.Cid) (*w3s.Status, error) {
	deals := []w3s.Deal{
		{
			Activation:        time.Now(),
			DealID:            1,
			DataModelSelector: "foo/bar",
		},
		{
			Activation:        time.Now().Add(time.Hour * 6),
			DealID:            2,
			DataModelSelector: "bar/foo",
		},
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
	jobs []UnfinihedJobs
}

func (m *mockCrdb) CreateJob(_ context.Context, cidStr string, pub string) error {
	cid, _ := cid.Decode(cidStr)
	m.jobs = append(m.jobs, UnfinihedJobs{
		Pub:       pub,
		Cid:       cid.Bytes(),
		Activated: time.Time{},
	})
	return nil
}

func (m *mockCrdb) UnfinishedJobs(_ context.Context) ([]UnfinihedJobs, error) {
	var t time.Time
	ufj := []UnfinihedJobs{}
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

func TestUploader(t *testing.T) {
	ctx := context.Background()
	mockGCS := new(mocks.GCS)

	// Mocking the returned values for the ParseEventData method
	mockGCS.On("ParseEvent").Return("mybucket", "myfile", nil)

	// Mocking the returned reader for the GetObjectReader method
	mockReadCloser := &MockReadCloser{Reader: bytes.NewReader(mockData())}
	mockGCS.On("GetObjectReader", ctx, "mybucket", "myfile").Return(mockReadCloser, nil)

	uploader := FileUploader{
		StorageClient: mockGCS,
		DealClient: &mockW3sClient{
			Files: []fs.File{},
		},
		DBClient: &mockCrdb{
			jobs: []UnfinihedJobs{},
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
	assert.Equal(t, "myfile", fStat.Name())
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
	assert.Equal(t, "myfile", db.jobs[0].Pub)

	cid, err := cid.Parse(db.jobs[0].Cid)
	assert.NoError(t, err)
	assert.Equal(t, getCIDFromBytes(mockData()).String(), cid.String())
}
