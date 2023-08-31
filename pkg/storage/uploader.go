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
	StorageClient GCSOps     // GCSClient is a GCSOps instance used to interact with GCS.
	DealClient    w3s.Client // W3SClient is a w3s.Client instance used to interact with W3S.
	DBClient      CrdbOps    // CrdbClient is a CrdbOps instance used to interact with CockroachDB.
}

// Upload downloads a file from GCS and uploads it to web3.storage.
func (u *FileUploader) Upload(ctx context.Context) error {
	bucket, fname, err := u.StorageClient.ParseEvent()
	if err != nil {
		return err
	}

	reader, err := u.StorageClient.GetObjectReader(
		ctx, bucket, fname)
	if err != nil {
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	fmt.Println("Read successful", bucket, fname)

	file := NewIntermediateFile(data, fname)
	cid, err := u.DealClient.Put(ctx, file)
	if err != nil {
		return err
	}

	// TODO: Add CID into cockroachdb
	fmt.Println("Upload successful :", cid)
	relName := "foobar" // comes from the CloudEvent
	err = u.DBClient.CreateDeal(ctx, cid.String(), relName)
	if err != nil {
		return err
	}

	fmt.Println("DB insert successful", relName)

	return nil
}
