package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ipfs/go-cid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/cockroachdb/cockroach-go/crdb"
)

func createDealTx(tx *sql.Tx, cidBytes []byte, relName string) error {
	_, err := tx.Exec(
		"Insert into deals (id, cid, relation) values(1, $1, $2)",
		cidBytes, relName)
	if err != nil {
		return errors.Wrap(err, "updating record")
	}

	return nil
}

type CrdbOps interface {
	CreateDeal(ctx context.Context, cidStr string, relationName string) error
}

type DBClient struct {
	// Initialize cockroachdb client to store metadata
	db *sql.DB
}

func NewDB(conn string) *DBClient {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &DBClient{
		db: db,
	}
}

func (db *DBClient) CreateDeal(ctx context.Context, cidStr string, relationName string) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	txopts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}
	err = crdb.ExecuteTx(ctx, db.db, txopts, func(tx *sql.Tx) error {
		return createDealTx(tx, cid.Bytes(), relationName)
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
