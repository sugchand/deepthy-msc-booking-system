package impl

import (
	"bookingSystem/roomInventory/pkg/env"
	"context"
	"database/sql"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresDB struct {
	rEnv     *env.RoomEnvValues
	dbhandle *bun.DB
	lock     sync.Mutex
}

// CreateDBContext will create a DB context for the DB operations
func (pd *PostgresDB) DB(ctx context.Context) *bun.DB {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	uname, pwd := pd.rEnv.DBUNameAndPwd()
	dsn := postgresDSN(pd.rEnv.DBRemote(), pd.rEnv.DBName(), uname, pwd)
	if pd.dbhandle == nil {
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		pd.dbhandle = bun.NewDB(sqldb, pgdialect.New())
	}

	return pd.dbhandle
}

func (pd *PostgresDB) Tx(ctx context.Context, db *bun.DB, opts *sql.TxOptions) (bun.Tx, error) {
	// we dont need to check error as we wanted to crash application when db
	// is not set in context.
	// No need of lock here.
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		log.WithError(err).Error("failed to begin new transaction")
		return bun.Tx{}, err
	}
	return tx, nil
}

func (pd *PostgresDB) TxCommit(ctx context.Context, tx bun.Tx) error {
	return tx.Commit()
}

func (pd *PostgresDB) TxRollBack(ctx context.Context, tx bun.Tx) error {
	return tx.Rollback()
}

func NewPostgresDB(uEnv *env.RoomEnvValues) *PostgresDB {
	return &PostgresDB{
		rEnv: uEnv,
	}
}

func postgresDSN(remote, dbName, uname, pwd string) string {
	// dsn := "postgres://postgres:@localhost:5432/test?sslmode=disable"
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	// TODO :: enable ssl to secure db access.
	return "postgres://" + uname + ":" + pwd + "@" + remote + "/" + dbName + "?sslmode=disable"
}
