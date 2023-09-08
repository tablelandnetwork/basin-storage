package storage

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"testing"
	"time"

	"crypto/sha256"

	"github.com/ipfs/go-cid"
	"github.com/tablelandnetwork/basin-storage/mocks"
	w3s "github.com/web3-storage/go-w3s-client"
	w3http "github.com/web3-storage/go-w3s-client/http"

	mh "github.com/multiformats/go-multihash"

	"github.com/stretchr/testify/assert"
)

// Mock interface for w3s.Client
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

func (m *mockW3sClient) Put(ctx context.Context, file fs.File, opts ...w3s.PutOption) (cid.Cid, error) {
	m.Files = append(m.Files, file)
	return getCIDFromBytes(mockData()), nil
}

func (m *mockW3sClient) Get(ctx context.Context, cid cid.Cid) (*w3http.Web3Response, error) {
	return nil, nil
}

func (m *mockW3sClient) PutCar(ctx context.Context, reader io.Reader) (cid.Cid, error) {
	return cid.Cid{}, nil
}
func (m *mockW3sClient) Status(ctx context.Context, cid cid.Cid) (*w3s.Status, error) {
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
func (m *mockW3sClient) List(ctx context.Context, opts ...w3s.ListOption) (*w3s.UploadIterator, error) {
	return nil, nil
}
func (m *mockW3sClient) Pin(context.Context, cid.Cid, ...w3s.PinOption) (*w3s.PinResponse, error) {
	return nil, nil
}

type MockReadCloser struct {
	*bytes.Reader
}

func (mrc *MockReadCloser) Close() error {
	return nil
}

type mockCrdb struct {
	db   map[string]string
	jobs []unfinihedJobs
}

func (m *mockCrdb) CreateJob(ctx context.Context, cidStr string, pub string) error {
	m.db[cidStr] = pub
	return nil
}

func (m *mockCrdb) UnfinishedJobs(ctx context.Context) ([]unfinihedJobs, error) {
	var t time.Time
	ufj := []unfinihedJobs{}
	for _, job := range m.jobs {
		if job.Activated == t {
			ufj = append(ufj, job)
		}
	}
	return ufj, nil
}

func (m *mockCrdb) UpdateJobStatus(ctx context.Context, cid []byte, activation time.Time) error {
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
			db: make(map[string]string),
		},
	}

	err := uploader.Upload(ctx)

	// Assert that mockGCS.GetObjectReader was called with the correct arguments
	mockGCS.AssertExpectations(t)
	assert.NoError(t, err)

	files := uploader.DealClient.(*mockW3sClient).Files
	assert.Equal(t, 1, len(files))

	fStat, _ := files[0].Stat()
	assert.Equal(t, "myfile", fStat.Name())
	assert.Equal(t, int64(11), fStat.Size())

	// Read the file contents that was send to (mocked) w3s
	buf := make([]byte, 11)
	files[0].Read(buf)
	files[0].Close()
	assert.Equal(t, mockData(), buf)

	// Assert that the CID was added to the database with the correct pub name
	pub := uploader.DBClient.(*mockCrdb).db[getCIDFromBytes(mockData()).String()]
	assert.Equal(t, "myfile", pub)
}
