package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"

	"github.com/tablelandnetwork/basin-storage/pkg/storage"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("Uploader", Uploader)
	functions.HTTP("StatusChecker", StatusChecker)
}

// Uploader is the CloudEvent function that is called by the Functions Framework.
// It is triggered by a CloudEvent that is published by the GCS bucket.
// The CloudEvent contains the name of the bucket and the name of the file.
// The file is downloaded from GCS and uploaded to web3.storage.
func Uploader(ctx context.Context, e event.Event) error {
	// Set a timeout of 60 minutes, thats the max time a function can run on GCP (gen2)
	// we want to ensure larger files can be uploaded
	cctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	// Read config from environment variables
	cfg := &storage.UploaderConfig{
		W3SToken: os.Getenv("WEB3STORAGE_TOKEN"),
		CrdbConn: os.Getenv("CRDB_CONN_STRING"),
	}

	// Initialize file uploader
	u, err := storage.NewFileUploader(cctx, e.Data(), cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize file uploader: %v", err)
	}

	// Upload file (from event) to web3.storage
	err = u.Upload(cctx)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// StatusChecker is the HTTP function that is called by the Functions Framework.
func StatusChecker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cfg := &storage.StatusCheckerConfig{
		W3SToken:         os.Getenv("WEB3STORAGE_TOKEN"),
		CrdbConn:         os.Getenv("CRDB_CONN_STRING"),
		PrivateKey:       "f4212951e2639746613d34dd912713c25df1be3d08e87680d4c9f147c81ed8b6",
		BackendURL:       "https://api.calibration.node.glif.io/rpc/v1",
		BasinStorageAddr: "0x4b1f4d8100e51afe644b189d77784dec225e0596",
	}

	sc, err := storage.NewStatusChecker(ctx, cfg)
	if err != nil {
		// TODO: return error
		panic(err)
	}

	err = sc.ProcessJobs(ctx)
	if err != nil {
		// TODO: return error
		panic(err)
	}

}
