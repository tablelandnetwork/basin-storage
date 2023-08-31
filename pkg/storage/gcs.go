package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

// GCSOps defines the interface for interacting with Google Cloud Storage (GCS).
type GCSOps interface {
	GetObjectReader(ctx context.Context, bName, oName string) (io.ReadCloser, error)
}

// GCSClient implements the GCSOps interface.
type GCSClient struct {
	Client *storage.Client
}

// GetObjectReader returns a reader for the specified object in the specified bucket.
func (r *GCSClient) GetObjectReader(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return r.Client.Bucket(bucketName).Object(objectName).NewReader(ctx)
}
