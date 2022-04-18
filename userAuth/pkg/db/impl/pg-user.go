package impl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

const userTableName = "users"

type UserTable struct {
	bun.BaseModel `bun:"table:users"` // table name users
	UserID        string              `bun:"user_id,pk"`
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
		UserID:    uuid.NewString(),
		UName:     username,
		PwdHash:   hash,
		Email:     email,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
	}

	db := dbHandle.DB(ctx)
	err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": userTableName,
				"user":     user.UName,
				"email":    user.Email,
			}).WithError(err).Error("failed to add new user to the table., cant execute transaction operation")
			return err
		}

		return nil
	})
	return err
}

// Get the User details.
func GetUser(ctx context.Context, dbHandle *PostgresDB, username, pwd string) (*UserTable, error) {
	db := dbHandle.DB(ctx)
	users := make([]UserTable, 0)
	// match both username and pwd hash
	err := db.NewSelect().Model(&users).Where("uname = ?", username).Scan(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
			"user":     username,
		}).WithError(err).Error("Cannot get the user from the db table.")
		return nil, err
	}
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

	// lets compare the password to make sure the user is provided with right password.
	err = bcrypt.CompareHashAndPassword([]byte(users[0].PwdHash), []byte(pwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
			"user":     username,
		}).WithError(err).Error("invalid password, cannot get user entry")
		return nil, err

	}
	return &users[0], nil
}

func DelUser(ctx context.Context, dbHandle *PostgresDB, username, pwd string) error {
	db := dbHandle.DB(ctx)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		users := make([]UserTable, 0)
		err := tx.NewSelect().Model(&users).Where("uname = ?", username).Scan(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": userTableName,
				"user":     username,
			}).WithError(err).Error("Failed to get the user entry from table, cannot delete!")
			return err
		}
		if len(users) == 0 {
			return nil // no user present., nothing to delete.
		}
		if len(users) != 1 {
			// hmm , that cannot be true, how does its possible to have more than one user have same credentials.
			err := errors.New("More than single user present in the system")
			log.WithFields(log.Fields{
				"db-table": userTableName,
				"user":     username,
			}).WithError(err).Error("Cannot delete more than one user from the table.")
			return err
		}
		err = bcrypt.CompareHashAndPassword([]byte(users[0].PwdHash), []byte(pwd))
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": userTableName,
				"user":     username,
			}).WithError(err).Error("invalid password, cannot delete user entry")
			return err
		}
		_, err = tx.NewDelete().Model((*UserTable)(nil)).Where("uname = ?", username).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": userTableName,
				"user":     username,
			}).WithError(err).Error("failed to delete the entry!")
			return err
		}
		return nil
	})
	return err
}

// reset password for a user.
func ResetPassword(ctx context.Context, dbHandle *PostgresDB, username, email, newPwd string) error {
	newHash, err := hashAndSalt([]byte(newPwd))
	if err != nil {
		log.WithFields(log.Fields{
			"db-table": userTableName,
		}).WithError(err).Error("failed to reset password in table., cannot calculate pwd hash")
		return err
	}
	db := dbHandle.DB(ctx)
	err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewUpdate().Model((*UserTable)(nil)).Set("pwdhash = ?", newHash).Where("uname = ?", username).Where("email = ?", email).Exec(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"db-table": userTableName,
			}).WithError(err).Error("failed to update password in the user table., cant execute DB transaction")
			return err
		}

		return nil
	})
	return err
}
