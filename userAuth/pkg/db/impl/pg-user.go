package impl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

const userTableName = "users"

type UserTable struct {
	bun.BaseModel `bun:"table:users"` // table name users
	ID            string              `bun:"id,pk,default:gen_random_uuid()"`
	UName         string              `bun:"uname,unique:uname"`
	PwdHash       string              `bun:"pwdhash"`
	Email         string              `bun:"email,unique:email"`
	IsAdmin       bool                `bun:"is_admin,default:FALSE"`
	CreatedAt     time.Time           `bun:"created_at"`
}

func CreateUserTable(ctx context.Context, dbHandle *PostgresDB) error {
	db := dbHandle.DB(ctx)
	_, err := db.NewCreateTable().Model((*UserTable)(nil)).
		IfNotExists().Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to create table.")
	}
	return err
}

func DropUserTable(ctx context.Context, dbHandle *PostgresDB) error {
	db := dbHandle.DB(ctx)
	_, err := db.NewDropTable().Model((*UserTable)(nil)).IfExists().Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to drop table.")
	}
	return err
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Error("Failed to generate hash+salt ", err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func NewUser(ctx context.Context, dbHandle *PostgresDB, username, pwd, email string, isAdmin bool) error {
	hash, err := hashAndSalt([]byte(pwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cannot calculate pwd hash")
	}
	user := &UserTable{
		UName:     username,
		PwdHash:   hash,
		Email:     email,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
	}

	db := dbHandle.DB(ctx)
	opts := &sql.TxOptions{}
	tx, err := dbHandle.Tx(ctx, db, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cant create a new DB transaction")
		return err
	}
	_, err = tx.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
			"user":     user.UName,
			"email":    user.Email,
		}).WithError(err).Error("failed to add new user to the table., cant execute transaction operation")
		return err
	}
	return tx.Commit()
}

// Get the User details.
func GetUser(ctx context.Context, dbHandle *PostgresDB, username, pwd string) (*UserTable, error) {
	hash, err := hashAndSalt([]byte(pwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cannot calculate pwd hash")
		return nil, err
	}
	db := dbHandle.DB(ctx)
	users := make([]UserTable, 0)
	// match both username and pwd hash
	db.NewSelect().Model(&users).Where("uname = ?", username).Where("pwdhash = ?", hash).Scan(ctx)
	// we expect single user for a specific username.
	if len(users) == 0 {
		return nil, nil // no user present.
	}
	if len(users) != 1 {
		// hmm , that cannot be true, how does its possible to have more than one user have same credentials.
		err := errors.New("More than single user present in the system")
		log.WithFields(log.Fields{
			"db-table": userTableName,
			"user":     username,
		}).WithError(err).Error("Cannot get the user from the db table.")
		return nil, err
	}
	return &users[0], nil
}

func DelUser(ctx context.Context, dbHandle *PostgresDB, username, pwd string) error {
	hash, err := hashAndSalt([]byte(pwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cannot calculate pwd hash")
		return err
	}
	db := dbHandle.DB(ctx)
	opts := &sql.TxOptions{}
	tx, err := dbHandle.Tx(ctx, db, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cant create a new DB transaction")
		return err
	}
	tx.NewDelete().Model((*UserTable)(nil)).Where("uname = ?", username).
		Where("pwdhash = ?", hash).Exec(ctx)
	return tx.Commit()
}

// reset password for a user.
func ResetPassword(ctx context.Context, dbHandle *PostgresDB, username, email, newPwd string) error {
	newHash, err := hashAndSalt([]byte(newPwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to add new user to the table., cannot calculate pwd hash")
		return err
	}
	db := dbHandle.DB(ctx)
	opts := &sql.TxOptions{}
	tx, err := dbHandle.Tx(ctx, db, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to update password in the user table., cant create a new DB transaction")
		return err
	}
	tx.NewUpdate().Model((*UserTable)(nil)).Set("pwdhash = ?", newHash).Where("uname = ?", username).Where("email = ?", email).Exec(ctx)
	return tx.Commit()
}
