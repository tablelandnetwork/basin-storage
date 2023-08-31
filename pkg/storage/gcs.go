package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

// GCSOps defines the interface for interacting with Google Cloud Storage (GCS).
type GCSOps interface {
	GetObjectReader(ctx context.Context, bName, oName string) (io.ReadCloser, error)
	ParseEvent() (string, string, error)
}

// GCSClient implements the GCSOps interface.
type GCSClient struct {
	Client    *storage.Client
	EventData []byte
}

// GetObjectReader returns a reader for the specified object in the specified bucket.
func (r *GCSClient) GetObjectReader(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return r.Client.Bucket(bucketName).Object(objectName).NewReader(ctx)
}

func (r *GCSClient) ParseEvent() (string, string, error) {
	var data storagedata.StorageObjectData
	if err := protojson.Unmarshal(r.EventData, &data); err != nil {
		return "", "", fmt.Errorf("protojson.Unmarshal: %w", err)
	}

	return data.GetBucket(), data.GetName(), nil
}
