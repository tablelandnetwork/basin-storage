package uploader

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"

	bstorage "github.com/tablelandnetwork/basin-storage/pkg/storage"
	w3s "github.com/web3-storage/go-w3s-client"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("Uploader", Uploader)
}

// Uploader is the CloudEvent function that is called by the Functions Framework.
// It is triggered by a CloudEvent that is published by the GCS bucket.
// The CloudEvent contains the name of the bucket and the name of the file.
// The file is downloaded from GCS and uploaded to web3.storage.
func Uploader(ctx context.Context, e event.Event) error {
	cctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	// Read w3s token and db conn string from environment variables
	web3StorageToken := os.Getenv("WEB3STORAGE_TOKEN")
	crdbConnStr := os.Getenv("CRDB_CONN_STRING")
	// connStr := "user=root password=\"\" dbname=defaultdb host=localhost port=26257 sslmode=disable"

	// Initialize GCS client to download file
	// bucket name and file name are passed in the CloudEvent
	storageClient, err := storage.NewClient(cctx)
	if err != nil {
		return fmt.Errorf("failed to initialize storage client: %v", err)
	}
	var data storagedata.StorageObjectData
	if err := protojson.Unmarshal(e.Data(), &data); err != nil {
		return fmt.Errorf("protojson.Unmarshal: %w", err)
	}

	// Initialize web3.storage client to upload file
	w3sOpts := []w3s.Option{
		w3s.WithToken(web3StorageToken),
		w3s.WithHTTPClient(
			&http.Client{
				Timeout: 0,
			},
		),
	}
	w3sClient, err := w3s.NewClient(w3sOpts...)
	if err != nil {
		return fmt.Errorf("failed to initialize web3.storage client: %v", err)
	}

	// Initialize cockroachdb client to store metadata
	db := bstorage.NewDB(crdbConnStr)

	u := &bstorage.FileUploader{
		Bucket:        data.GetBucket(),
		Filename:      data.GetName(),
		StorageClient: &bstorage.GCSClient{Client: storageClient},
		DealClient:    w3sClient,
		DBClient:      db,
	}
	err = u.Upload(cctx)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	return nil
}
