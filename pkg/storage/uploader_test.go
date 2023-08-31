package storage

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"testing"

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

func getCIDFromMockData() cid.Cid {
	hashedData := sha256.Sum256(mockData())
	multihash, _ := mh.Encode(hashedData[:], mh.SHA2_256)
	cidV1 := cid.NewCidV1(cid.Raw, multihash)
	return cidV1
}

func (m *mockW3sClient) Put(ctx context.Context, file fs.File, opts ...w3s.PutOption) (cid.Cid, error) {
	m.Files = append(m.Files, file)
	return getCIDFromMockData(), nil
}

func (m *mockW3sClient) Get(ctx context.Context, cid cid.Cid) (*w3http.Web3Response, error) {
	return nil, nil
}

func (m *mockW3sClient) PutCar(ctx context.Context, reader io.Reader) (cid.Cid, error) {
	return cid.Cid{}, nil
}
func (m *mockW3sClient) Status(ctx context.Context, cid cid.Cid) (*w3s.Status, error) {
	return nil, nil
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

type MockCrdb struct {
	db map[string]string
}

func (m *MockCrdb) CreateDeal(ctx context.Context, cidStr string, relationName string) error {
	m.db[cidStr] = relationName
	return nil
}

func TestMyFunction(t *testing.T) {
	ctx := context.Background()
	mockGCS := new(mocks.GCSOps)

	// Mocking the returned reader for the GetObjectReader method
	mockReadCloser := &MockReadCloser{Reader: bytes.NewReader(mockData())}
	mockGCS.On("GetObjectReader", ctx, "mybucket", "myfile").Return(mockReadCloser, nil)

	uploader := FileUploader{
		Bucket:        "mybucket",
		Filename:      "myfile",
		StorageClient: mockGCS,
		DealClient: &mockW3sClient{
			Files: []fs.File{},
		},
		DBClient: &MockCrdb{
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

	// Assert that the CID was added to the database with the correct relation name
	assert.Equal(t, "foobar", uploader.DBClient.(*MockCrdb).db[getCIDFromMockData().String()])
}
