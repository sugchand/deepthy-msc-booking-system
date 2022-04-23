package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type DBHandle interface {
	// function for the DB handle
	DB(ctx context.Context) *bun.DB
	Tx(ctx context.Context, db *bun.DB, opts *sql.TxOptions) (bun.Tx, error)
	TxCommit(ctx context.Context, tx bun.Tx) error
	TxRollBack(ctx context.Context, tx bun.Tx) error
}
