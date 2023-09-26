package tests

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "github.com/tablelandnetwork/basin-storage"
)

func buildCheckerRequest(t *testing.T) *http.Request {
	urlStr := fmt.Sprintf("http://localhost:%s", functionsPort)
	data := url.Values{}
	data.Set("simulated", "true")
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestChecker(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	w3sToken := os.Getenv("WEB3STORAGE_TOKEN")
	dbHost := os.Getenv("CRDB_HOST")
	pk := os.Getenv("PRIVATE_KEY")
	chainIDStr := os.Getenv("CHAIN_ID")
	backendURL := "https://api.calibration.node.glif.io/rpc/v1"
	basinStorageAddr := "0x4b1f4d8100e51afe644b189d77784dec225e0596"
	crdbConn := fmt.Sprintf(
		"postgresql://root@%s/basin_test?sslmode=disable",
		dbHost)

	// setup db for testing
	db, err := sql.Open("postgres", crdbConn)
	require.NoError(t, err)

	// setup initial database state
	SetupDB(t, db)

	// insert a processed job
	cid := insertProcessedJob(t, db)
	defer func() {
		_, err := db.Exec("DROP DATABASE IF EXISTS basin_test")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	}()

	// Create a pub in the smart contract if it doesn't exist
	createPub(t, pk, chainIDStr, backendURL, basinStorageAddr)

	// start the cloud function
	go func() {
		err := funcframework.RegisterHTTPFunctionContext(
			context.Background(),
			"/",
			handler.StatusChecker,
		)
		require.NoError(t, err)
		require.NoError(t, os.Setenv("W3S_TOKEN", w3sToken))
		require.NoError(t, os.Setenv("CRDB_CONN_STRING", crdbConn))
		require.NoError(t, os.Setenv("PRIVATE_KEY", pk))
		require.NoError(t, os.Setenv("CHAIN_ID", chainIDStr))
		require.NoError(t, funcframework.Start(functionsPort))
	}()

	time.Sleep(2 * time.Second)

	req := buildCheckerRequest(t)
	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	query := `
		SELECT namespaces.name, jobs.cid, jobs.relation, jobs.activated
		FROM namespaces, jobs
		WHERE namespaces.id = jobs.ns_id
		AND jobs.cid = $1
	`
	rows, err := db.Query(query, cid.Bytes())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, rows.Close())
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
	assert.Equal(t, "esfbmltndstj", results[0].nsName)
	assert.Equal(t, "ksvraapqfiyf", results[0].relName)
	assert.NotNil(t, results[0].cid)
	value, err := results[0].activated.Value()
	require.NoError(t, err)
	assert.Equal(t, "2023-09-26T08:09:30Z", value)
}
