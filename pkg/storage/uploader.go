package storage

import (
	"context"
	"fmt"
	"io"

	w3s "github.com/web3-storage/go-w3s-client"
)

// FileUploader dowload a file form GCS and uploads to web3.storage.
type FileUploader struct {
	Bucket        string     // SourceBucket is the name of the GCS bucket where the file will be uploaded from.
	Filename      string     // Filename is the name of the file to be uploaded.
	StorageClient GCS        // GCSClient is a GCSOps instance used to interact with GCS.
	DealClient    w3s.Client // W3SClient is a w3s.Client instance used to interact with W3S.
	DBClient      Crdb       // CrdbClient is a CrdbOps instance used to interact with CockroachDB.
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
	defer reader.Close()

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

	err = u.DBClient.CreateDeal(ctx, cid.String(), fname)
	if err != nil {
		return fmt.Errorf("failed to create deal: %v", err)
	}

	fmt.Println("DB insert successful", fname)

	return nil
}
