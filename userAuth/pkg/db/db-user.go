package db

import (
	"bookingSystem/userAuth/pkg/db/impl"
	"bookingSystem/userAuth/pkg/env"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type userToken struct {
	token   string
	expiry  time.Time
	isAdmin bool
}

type UserTableHandle struct {
	pgDBHandle *impl.PostgresDB
	lock       sync.Mutex // lock to protect the following session map
	// we use username as a key here to keep track of access tokens
	accessTokenMap map[string]*userToken
	envValues      *env.UserEnvValues
}

func (uth *UserTableHandle) createNewSession(uname string, isAdmin bool) (string, time.Duration) {
	tok := &userToken{
		token:   uuid.New().String(),
		expiry:  time.Now().Add(uth.envValues.TokenValidity()),
		isAdmin: isAdmin,
	}
	// add new session entry.
	uth.accessTokenMap[uname] = tok
	return tok.token, uth.envValues.TokenValidity()
}

func (uth *UserTableHandle) NewUser(ctx context.Context, username, pwd, email string, isAdmin bool) error {
	return impl.NewUser(ctx, uth.pgDBHandle, username, pwd, email, isAdmin)
}

func (uth *UserTableHandle) DeleteUser(ctx context.Context, username, pwd string) error {
	return impl.DelUser(ctx, uth.pgDBHandle, username, pwd)
}

func (uth *UserTableHandle) ResetPassword(ctx context.Context, uname, email, newPwd string) error {
	return impl.ResetPassword(ctx, uth.pgDBHandle, uname, email, newPwd)
}

func (uth *UserTableHandle) GetUserToken(ctx context.Context, uname, pwd string) (string, time.Duration, bool, error) {
	uth.lock.Lock()
	defer uth.lock.Unlock()
	// check if a session already present for user before reading from DB.
	if token, ok := uth.accessTokenMap[uname]; ok {
		timeLeft := time.Until(token.expiry)
		if timeLeft >= 0 {
			// so the session is still valid
			return token.token, timeLeft, token.isAdmin, nil
		}
		// So now we have the session expired situation, lets delete the entry from the map
		delete(uth.accessTokenMap, uname)
	}
	// lets check if user present in DB table.
	userRow, err := impl.GetUser(ctx, uth.pgDBHandle, uname, pwd)
	if err != nil {
		return "", 0, false, err
	}
	if userRow == nil {
		// user is not present and no session
		return "", 0, false, errors.New("user account not present")
	}
	// finally add it to session map.
	sessionToken, validity := uth.createNewSession(uname, userRow.IsAdmin)
	return sessionToken, validity, userRow.IsAdmin, nil

}

// For now we use postgress as DB backend. its possible to replace it with
// another backend by changing this file with new handler for the DB
func NewDBUserTableHandle(ctx context.Context, userEnvValues *env.UserEnvValues) (*UserTableHandle, error) {
	pgHandle := impl.NewPostgresDB(userEnvValues)
	pgHandle.DB(ctx) // initialize the DB Handle
	// lets create a table as well for the user if not exists.
	err := impl.CreateUserTable(ctx, pgHandle)
	if err != nil {
		// so failed to create a table and cant procced without the user table.
		return nil, errors.New("failed to create user table in postgres!.")
	}
	return &UserTableHandle{
		pgDBHandle:     pgHandle,
		accessTokenMap: make(map[string]*userToken),
		envValues:      userEnvValues,
	}, nil
}
