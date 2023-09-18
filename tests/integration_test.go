package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Blank-import libpq package for SQL.
	_ "github.com/lib/pq"
	handler "github.com/tablelandnetwork/basin-storage"
)

func uploadRandomBytesToGCS(t *testing.T, data []byte, bucketName, objectName string) {
	ctx := context.Background()

	// Create a client
	client, err := storage.NewClient(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close())
	}()

	// Get a handle to the bucket and the object
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	// Create a writer and write the random bytes
	wc := object.NewWriter(ctx)
	_, err = wc.Write(data)
	require.NoError(t, err)
	require.NoError(t, wc.Close())
}

func deleteObjectFromGCS(t *testing.T, bucketName, objectName string) {
	ctx := context.Background()

	// Create a client
	client, err := storage.NewClient(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close())
	}()

	// Get a handle to the bucket and the object
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	// Delete the object
	err = object.Delete(ctx)
	require.NoError(t, err)
}

func SetupDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS basin_test")
	require.NoError(t, err)

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS namespaces
	(
		id     BIGSERIAL   PRIMARY KEY,
		name   VARCHAR(32) UNIQUE NOT NULL,
		owner  BYTEA NOT NULL,
		last_export TIMESTAMPTZ
	)`)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO namespaces (name, owner) VALUES ('test_name', 'test_owner')")
	require.NoError(t, err)

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS jobs
		(
			id BIGSERIAL PRIMARY KEY,
			ns_id BIGINT,
			cid BYTEA NOT NULL,
			relation TEXT NOT NULL,
			activated TIMESTAMP,
			CONSTRAINT fk_namespace
			FOREIGN KEY(ns_id)
			REFERENCES namespaces(id)
		)`)
	require.NoError(t, err)
}

func buildUplaodRequest(t *testing.T, bucketName, objectName string) *http.Request {
	url := "http://localhost:8293"
	postData := fmt.Sprintf(
		`{
			"name": "%s",
			"bucket": "%s",
			"contentType": "application/json",
			"metageneration": "1",
			"timeCreated": "2020-04-23T07:38:57.230Z",
			"updated": "2020-04-23T07:38:57.230Z"
		}`,
		objectName,
		bucketName,
	)

	data := []byte(postData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	require.NoError(t, err)

	// Set headers
	source := fmt.Sprintf(
		"//storage.googleapis.com/projects/_/buckets/%s",
		bucketName)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ce-id", "1234567890")
	req.Header.Set("ce-specversion", "1.0")
	req.Header.Set("ce-type", "google.cloud.storage.object.v1.finalized")
	req.Header.Set("ce-time", "2020-08-08T00:11:44.895529672Z")
	req.Header.Set("ce-source", source)

	return req
}

func TestUploader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	w3sToken := os.Getenv("WEB3STORAGE_TOKEN")
	dbHost := os.Getenv("CRDB_HOST")
	crdbConn := fmt.Sprintf(
		"postgresql://root@%s/basin_test?sslmode=disable",
		dbHost)

	// setup db for testing
	db, err := sql.Open("postgres", crdbConn)
	require.NoError(t, err)
	SetupDB(t, db)
	defer func() {
		_, err := db.Exec("DROP DATABASE IF EXISTS basin_test")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	}()

	// start the cloud function
	go func() {
		err := funcframework.RegisterCloudEventFunctionContext(
			context.Background(),
			"/",
			handler.Uploader,
		)
		require.NoError(t, err)
		require.NoError(t, os.Setenv("W3S_TOKEN", w3sToken))
		require.NoError(t, os.Setenv("CRDB_CONN_STRING", crdbConn))
		require.NoError(t, funcframework.Start("8293"))
	}()

	// Upload random bytes to GCS for testing
	bucketName := "tableland-entrypoint"
	objectName := "01234-012345-1-2-00000000-basin_storage.test_name.data6-test.parquet"
	size := 1 * 1024 * 1024 // 1MB
	data := make([]byte, size)
	_, err = rand.Read(data)
	require.NoError(t, err)
	uploadRandomBytesToGCS(t, data, bucketName, objectName)
	defer deleteObjectFromGCS(t, bucketName, objectName)

	// Wait for for test file to be uploaded to GCS
	time.Sleep(3 * time.Second)

	// Trigger the cloud function
	req := buildUplaodRequest(t, bucketName, objectName)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	// For example, check if the response status is 200 OK
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the database to check if the job was created
	query := `
		SELECT namespaces.name, jobs.cid, jobs.relation, jobs.activated
		FROM namespaces, jobs
		WHERE namespaces.id = jobs.ns_id and activated is NULL
	`
	rows, err := db.Query(query)
	require.NoError(t, err)
	defer func() {
		if err := rows.Close(); err != nil {
			require.NoError(t, err)
		}
	}()

	type result struct {
		nsName    string
		relName   string
		cid       []byte
		activated sql.NullString
	}

	var results []result
	for rows.Next() {
		var cid []byte
		var nsName string
		var relation string
		var activated sql.NullString
		if err := rows.Scan(&nsName, &cid, &relation, &activated); err != nil {
			require.NoError(t, err)
		}
		results = append(results, result{
			nsName:    nsName,
			relName:   relation,
			cid:       cid,
			activated: activated,
		})
	}
	defer func() {
		require.NoError(t, rows.Close())
	}()

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "test_name", results[0].nsName)
	assert.Equal(t, "data6", results[0].relName)
	assert.NotNil(t, results[0].cid)
	assert.False(t, results[0].activated.Valid)
}
