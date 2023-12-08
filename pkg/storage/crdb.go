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

func createJobTx(
	tx *sql.Tx,
	cidBytes []byte,
	pub Pub,
	timestamp *int64,
	fname string,
	cacheDuration int64,
	sign string,
	hash string,
) error {
	row := tx.QueryRow(
		"SELECT id FROM namespaces WHERE name = $1", pub.Namespace)
	var nsID int
	if err := row.Scan(&nsID); err != nil {
		return errors.Wrap(err, "error while querying namespace")
	}

	expiresAt, cachePath := &sql.NullTime{}, &sql.NullString{}
	if cacheDuration > 0 {
		_ = expiresAt.Scan(time.Now().Add(time.Minute * time.Duration(cacheDuration)).UTC())
		_ = cachePath.Scan(fname)
	}

	_, err := tx.Exec(
		"insert into jobs (ns_id, cid, relation, timestamp, cache_path, expires_at, sign, hash) values($1, $2, $3, $4, $5, $6, $7, $8)",
		nsID, cidBytes, pub.Relation, timestamp, sign, hash)
	if err != nil {
		return errors.Wrap(err, "updating record")
	}

	return nil
}

// Crdb is an interface that defines the methods to interact with CockroachDB.
type Crdb interface {
	CreateJob(ctx context.Context, cidStr string, fileName string, timestamp *int64, cacheDuration int64, sign, hash string) error
	UnfinishedJobs(ctx context.Context) ([]UnfinishedJob, error)
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

func extractPub(filename string) (Pub, error) {
	fmt.Println("Extracting pub name from filename: ", filename)
	parts := strings.Split(filename, "/")
	if len(parts) < 3 {
		return Pub{}, fmt.Errorf("invalid filename")
	}

	// parts[0] is the namespace
	// parts[1] is the table/relation name
	return Pub{
		Namespace: parts[0],
		Relation:  parts[1],
	}, nil
}

// CreateJob creates a new job in the DB.
func (db *DBClient) CreateJob(
	ctx context.Context,
	cidStr string,
	fname string,
	timestamp *int64,
	cacheDuration int64,
	sign string,
	hash string,
) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return fmt.Errorf("failed to decode cid: %v", err)
	}

	txopts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	// Extract the schema and table name from the file name.
	pub, err := extractPub(fname)
	if err != nil {
		return fmt.Errorf("failed to extract table name: %v", err)
	}

	err = crdb.ExecuteTx(ctx, db.DB, txopts, func(tx *sql.Tx) error {
		return createJobTx(tx, cid.Bytes(), pub, timestamp, fname, cacheDuration, sign, hash)
	})
	if err != nil {
		return fmt.Errorf("failed to create new job: %v", err)
	}

	return nil
}

// UnfinishedJob represents a job in db that
// tracks deals on Filecoin.
type UnfinishedJob struct {
	Pub       Pub
	Cid       []byte
	Activated time.Time
	Timestamp *int64
	CachePath string
	ExpiresAt time.Time
}

// UnfinishedJobs returns all currently unfinished jobs in the db.
func (db *DBClient) UnfinishedJobs(ctx context.Context) ([]UnfinishedJob, error) {
	query := `
		SELECT namespaces.name, jobs.cid, jobs.relation, jobs.timestamp
		FROM namespaces, jobs
		WHERE namespaces.id = jobs.ns_id and activated is NULL
	`
	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query unfinished jobs: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("error when closing crdb connection: %v", err)
		}
	}()

	var result []UnfinishedJob
	for rows.Next() {
		var cid []byte
		var nsName string
		var relation string
		var timestamp sql.NullInt64
		if err := rows.Scan(&nsName, &cid, &relation, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		var t *int64
		if timestamp.Valid {
			t = &timestamp.Int64
		} else {
			t = nil
		}
		result = append(result, UnfinishedJob{
			Pub: Pub{
				Namespace: nsName,
				Relation:  relation,
			},
			Cid:       cid,
			Timestamp: t,
		})
	}

	return result, nil
}

// UpdateJobStatus updates the job status in the DB.
func (db *DBClient) UpdateJobStatus(ctx context.Context, cid []byte, activation time.Time) error {
	_, err := db.DB.ExecContext(ctx,
		"UPDATE jobs SET activated = $1 WHERE cid = $2",
		activation, cid,
	)
	if err != nil {
		return fmt.Errorf("failed to update job status: %v", err)
	}

	return nil
}
