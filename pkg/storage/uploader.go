package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	w3s "github.com/web3-storage/go-w3s-client"
)

// FileUploader download a file from GCS and uploads to web3.storage.
type FileUploader struct {
	StorageClient GCS        // StorageClient is a GCS instance used to interact with GCS.
	DealClient    w3s.Client // DealClient is a w3s.Client instance used to interact with W3S.
	DBClient      Crdb       // DBClient is a Crdb instance used to interact with CockroachDB.
}

// UploaderConfig defines the configuration for a FileUploader.
type UploaderConfig struct {
	W3SToken string
	CrdbConn string
}

// NewFileUploader creates a new FileUploader.
func NewFileUploader(ctx context.Context, eventData []byte, cfg *UploaderConfig) (*FileUploader, error) {
	// Initialize GCS client to download file
	// bucket name and file name are passed in the CloudEvent
	storageClient, err := NewGCSClient(ctx, eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage client: %v", err)
	}

	// Initialize web3.storage client to upload file
	w3sOpts := []w3s.Option{
		w3s.WithToken(cfg.W3SToken),
		w3s.WithHTTPClient(
			&http.Client{
				Timeout: 0, // no timeout
			},
		),
	}
	w3sClient, err := w3s.NewClient(w3sOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize web3.storage client: %v", err)
	}

	// Initialize cockroachdb client to store metadata
	dbClient, err := NewDB(cfg.CrdbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cockroachdb client: %v", err)
	}

	u := &FileUploader{
		StorageClient: storageClient,
		DealClient:    w3sClient,
		DBClient:      dbClient,
	}

	return u, nil
}

// Upload downloads a file from GCS and uploads it to web3.storage.
func (u *FileUploader) Upload(ctx context.Context) error {
	bucket, fname, err := u.StorageClient.ParseEvent()
	if err != nil {
		return fmt.Errorf("failed to parse event: %v", err)
	}

	reader, err := u.StorageClient.GetObjectReader(
		ctx, bucket, fname)
	if err != nil {
		return fmt.Errorf("failed to get object reader: %v", err)
	}

	metadata, err := u.StorageClient.GetObjectMetadata(ctx, bucket, fname)
	if err != nil {
		return fmt.Errorf("failed to get object metadata: %v", err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Fatalf("error when closing cloud storage reader: %v", err)
		}
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read object: %v", err)
	}

	fmt.Println("Read successful", bucket, fname)

	file := NewIntermediateFile(data, fname)
	cid, err := u.DealClient.Put(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	fmt.Println("Upload successful :", cid)

	var timestamp *int64
	if _, ok := metadata["timestamp"]; !ok {
		fmt.Println("timestamp is missing", fname)
	} else {
		t, err := strconv.ParseInt(metadata["timestamp"], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		timestamp = &t
	}

	var cacheDutation int64
	if _, ok := metadata["cache_duration"]; ok {
		duration, err := strconv.ParseInt(metadata["cache_duration"], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		cacheDutation = duration
	}

	sign, ok := metadata["signature"]
	if !ok {
		return fmt.Errorf("signature is missing")
	}

	hash, ok := metadata["hash"]
	if !ok {
		return fmt.Errorf("hash is missing")
	}

	err = u.DBClient.CreateJob(ctx, cid.String(), fname, timestamp, cacheDutation, sign, hash)
	if err != nil {
		return err
	}

	fmt.Println("DB insert successful", fname)

	return nil
}
