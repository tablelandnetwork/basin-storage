package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	// Blank-import libpq package for SQL.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach-go/crdb"
)

func createJobTx(tx *sql.Tx, cidBytes []byte, pub Pub) error {
	row := tx.QueryRow(
		"SELECT id FROM namespaces WHERE name = $1", pub.Namespace)
	var nsID int
	if err := row.Scan(&nsID); err != nil {
		return errors.Wrap(err, "error while querying namespace")
	}

	_, err := tx.Exec(
		"Insert into deals (ns_id, cid, relation) values($1, $2, $3)",
		nsID, cidBytes, pub.Relation)
	if err != nil {
		return errors.Wrap(err, "updating record")
	}

	return nil
}

// Crdb is an interface that defines the methods to interact with CockroachDB.
type Crdb interface {
	CreateJob(ctx context.Context, cidStr string, fileName string) error
	UnfinishedJobs(ctx context.Context) ([]UnfinihedJobs, error)
	UpdateJobStatus(ctx context.Context, cid []byte, activation time.Time) error
}

// DBClient is a Crdb implementation.
type DBClient struct {
	// Initialize cockroachdb client to store metadata
	DB *sql.DB
}

// NewDB creates a new DBClient.
func NewDB(conn string) (*DBClient, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		DB: db,
	}, nil
}

// Pub represents a ns and table/relation name.
type Pub struct {
	Namespace string
	Relation  string
}

func (db *DBClient) extractPubName(filename string) (Pub, error) {
	filenameParts := strings.Split(filename, "-")
	if len(filenameParts) < 2 {
		return Pub{}, fmt.Errorf("invalid filename")
	}

	// parts[0] is the database name, which we don't need
	// parts[1:len(parts)-1] is the represents namespace
	// parts[len(parts)-1] is the table/relation name
	parts := strings.Split(filenameParts[len(filenameParts)-2], ".")
	if len(parts) < 3 {
		return Pub{}, fmt.Errorf("invalid schema or table name")
	}

	partsLen := len(parts)
	Pub := Pub{
		Namespace: strings.Join(parts[1:partsLen-1], "."),
		Relation:  parts[partsLen-1],
	}
	return Pub, nil
}

// CreateJob creates a new job in the DB.
func (db *DBClient) CreateJob(ctx context.Context, cidStr string, fname string) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return fmt.Errorf("failed to decode cid: %v", err)
	}

	txopts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	// Extract the schema and table name from the file name.
	pub, err := db.extractPubName(fname)
	if err != nil {
		return fmt.Errorf("failed to extract table name: %v", err)
	}

	err = crdb.ExecuteTx(ctx, db.DB, txopts, func(tx *sql.Tx) error {
		return createJobTx(tx, cid.Bytes(), pub)
	})
	if err != nil {
		return fmt.Errorf("failed to create new job: %v", err)
	}

	return nil
}

// UnfinihedJobs represents an unfinished job.
type UnfinihedJobs struct {
	Pub       string
	Cid       []byte
	Activated time.Time
}

// UnfinishedJobs returns all unfinished jobs.
func (db *DBClient) UnfinishedJobs(ctx context.Context) ([]UnfinihedJobs, error) {
	rows, err := db.DB.QueryContext(ctx,
		"SELECT namespaces.name, deals.cid, deals.relation FROM namespaces, deals WHERE namespaces.id = deals.ns_id and activated is NULL")
	if err != nil {
		return nil, fmt.Errorf("failed to query unfinished jobs: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("error when closing crdb connection: %v", err)
		}
	}()

	var result []UnfinihedJobs
	for rows.Next() {
		var cid []byte
		var nsName string
		var relation string
		if err := rows.Scan(&nsName, &cid, &relation); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result = append(result, UnfinihedJobs{
			Pub: fmt.Sprintf("%s.%s", nsName, relation),
			Cid: cid,
		})
	}

	return result, nil
}

// UpdateJobStatus updates the job status in the DB.
func (db *DBClient) UpdateJobStatus(ctx context.Context, cid []byte, activation time.Time) error {
	_, err := db.DB.ExecContext(ctx,
		"UPDATE deals SET activated = $1 WHERE cid = $2",
		activation, cid,
	)
	if err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	return nil
}
