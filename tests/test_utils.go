package tests

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"github.com/tablelandnetwork/basin-storage/pkg/ethereum"
	"github.com/textileio/go-tableland/pkg/wallet"
	"golang.org/x/crypto/sha3"
)

const functionsPort = "8293"

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

	_, err = db.Exec("INSERT INTO namespaces (name, owner) VALUES ('esfbmltndstj', 'test_owner')")
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

func insertProcessedJob(t *testing.T, db *sql.DB) cid.Cid {
	// already processed job with 2 deals
	cidStr := "bafybeie62jygi4bi5m3yd43ttppkn5velvck5yutckggp6iwdmkkowyhxy"
	cid, err := cid.Decode(cidStr)
	require.NoError(t, err)

	_, err = db.Exec(
		`INSERT INTO jobs (ns_id, cid, relation)
		VALUES (
			(SELECT id FROM namespaces WHERE name = 'esfbmltndstj'),
			$1,
			'ksvraapqfiyf'			
		)`, cid.Bytes())
	require.NoError(t, err)

	return cid
}

func createPub(t *testing.T, pk, chainIDStr, backendURL, basinStorageAddr string) {
	wallet, err := wallet.NewWallet(pk)
	require.NoError(t, err)

	backend, err := ethclient.DialContext(context.Background(), backendURL)
	require.NoError(t, err)

	addr, err := common.NewMixedcaseAddressFromString(basinStorageAddr)
	require.NoError(t, err)

	chainID, err := strconv.ParseUint(chainIDStr, int(10), 64)
	require.NoError(t, err)

	contract, err := ethereum.NewContract(
		addr.Address(),
		backend,
	)
	require.NoError(t, err)

	txOpts, err := bind.NewKeyedTransactorWithChainID(
		wallet.PrivateKey(),
		big.NewInt(int64(chainID)),
	)
	require.NoError(t, err)

	res, err := contract.CreatePub(
		txOpts,
		wallet.Address(),
		"esfbmltndstj.ksvraapqfiyf",
	)
	if err != nil {
		hasher := sha3.NewLegacyKeccak256()
		hasher.Write([]byte("PubAlreadyExists(string)"))
		hash := hasher.Sum(nil)
		if !strings.Contains(err.Error(), fmt.Sprintf("0x%x", hash[:4])) {
			// ignores error if pub already exists error
			require.NoError(t, err)
		}
	}
	if res != nil {
		// wait for the transaction to be included in a block
		time.Sleep(150 * time.Second)
	}
}
