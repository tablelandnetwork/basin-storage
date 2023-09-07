package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach-go/crdb"
)

func createJobTx(tx *sql.Tx, cidBytes []byte, relName string) error {
	_, err := tx.Exec(
		"Insert into deals (cid, relation) values($1, $2)",
		cidBytes, relName)
	if err != nil {
		return errors.Wrap(err, "updating record")
	}

	return nil
}

type Crdb interface {
	CreateJob(ctx context.Context, cidStr string, relationName string) error
	UnfinishedJobs(ctx context.Context) ([][]byte, error)
	UpdateJobStatus(ctx context.Context, cidStr string, activation time.Time) error
}

type DBClient struct {
	// Initialize cockroachdb client to store metadata
	DB *sql.DB
}

func NewDB(conn string) (*DBClient, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		DB: db,
	}, nil
}

func (db *DBClient) extractTblName(filename string) (string, error) {
	parts := strings.Split(filename, "-")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid filename")
	}
	return parts[len(parts)-2], nil
}

func (db *DBClient) CreateJob(ctx context.Context, cidStr string, fileName string) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return fmt.Errorf("failed to decode cid: %v", err)
	}

	txopts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tblName, err := db.extractTblName(fileName)
	if err != nil {
		return fmt.Errorf("failed to extract table name: %v", err)
	}

	err = crdb.ExecuteTx(ctx, db.DB, txopts, func(tx *sql.Tx) error {
		return createJobTx(tx, cid.Bytes(), tblName)
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}

	return nil
}

func (db *DBClient) UnfinishedJobs(ctx context.Context) ([][]byte, error) {
	rows, err := db.DB.Query("SELECT cid FROM deals WHERE activated is NULL")
	if err != nil {
		return nil, fmt.Errorf("failed to query deals: %v", err)
	}
	defer rows.Close()

	var cids [][]byte
	for rows.Next() {
		var cid []byte
		if err := rows.Scan(&cid); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		cids = append(cids, cid)
	}

	return cids, nil
}

func (db *DBClient) UpdateJobStatus(ctx context.Context, cidStr string, activation time.Time) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return fmt.Errorf("failed to decode cid: %v", err)
	}

	_, err = db.DB.Exec(
		"UPDATE deals SET activation = $1 WHERE cid = $2",
		activation, cid.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	return nil
}
