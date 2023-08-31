package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

type Crdb interface {
	CreateDeal(ctx context.Context, cidStr string, relationName string) error
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

func (db *DBClient) CreateDeal(ctx context.Context, cidStr string, fileName string) error {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	txopts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tblName, err := db.extractTblName(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = crdb.ExecuteTx(ctx, db.DB, txopts, func(tx *sql.Tx) error {
		return createDealTx(tx, cid.Bytes(), tblName)
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
