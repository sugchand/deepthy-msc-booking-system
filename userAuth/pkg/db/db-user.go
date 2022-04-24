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

type activeTokenData struct {
	uname   string
	expiry  time.Time
	isAdmin bool
}

type UserTableHandle struct {
	pgDBHandle *impl.PostgresDB
	lock       sync.Mutex // lock to protect the following session map
	// we use username+pwd as a key here to keep track of logins.
	accessTokenMap map[string]*userToken
	// use token as key to keep track of active tokens and their relevant uname.
	activeTokens map[string]*activeTokenData
	envValues    *env.UserEnvValues
}

func (uth *UserTableHandle) createNewSession(key string, isAdmin bool) *userToken {
	tok := &userToken{
		token:   uuid.New().String(),
		expiry:  time.Now().Add(uth.envValues.TokenValidity()),
		isAdmin: isAdmin,
	}
	// add new session entry.
	uth.accessTokenMap[key] = tok
	return tok
}

func (uth *UserTableHandle) createActiveTokenEntry(token, uname string, expiry time.Time, isAdmin bool) {
	tok := &activeTokenData{
		uname:   uname,
		expiry:  expiry,
		isAdmin: isAdmin,
	}
	uth.activeTokens[token] = tok
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

// a simple function to keep track of username and pwd.
func unamePwdKey(uname, pwd string) string {
	return uname + pwd
}

func (uth *UserTableHandle) GetUserToken(ctx context.Context, uname, pwd string) (string, string, time.Duration, bool, error) {
	// TODO :: for now we use single token for a user across all his sessions in different machines and tabs.
	// Its fine for now as we let them use single session everywhere. Remove that logic and create token per session.
	key := unamePwdKey(uname, pwd)
	uth.lock.Lock()
	defer uth.lock.Unlock()
	// check if a session already present for user before reading from DB.
	if token, ok := uth.accessTokenMap[key]; ok {
		timeLeft := time.Until(token.expiry)
		if timeLeft >= 0 {
			// so the session is still valid
			return uname, token.token, timeLeft, token.isAdmin, nil
		}
		// So now we have the session expired situation, lets delete the entry from the map
		delete(uth.accessTokenMap, uname)
	}
	// lets check if user present in DB table.
	userRow, err := impl.GetUser(ctx, uth.pgDBHandle, uname, pwd)
	if err != nil {
		return "", "", 0, false, err
	}
	if userRow == nil {
		// user is not present and no session
		return "", "", 0, false, errors.New("user account not present")
	}
	// finally add it to session map.
	tok := uth.createNewSession(key, userRow.IsAdmin)
	uth.createActiveTokenEntry(tok.token, uname, tok.expiry, tok.isAdmin)
	return uname, tok.token, time.Until(tok.expiry), tok.isAdmin, nil

}

func (uth *UserTableHandle) TokenValid(ctx context.Context, token, uname string) (string, time.Duration, bool, error) {
	uth.lock.Lock()
	defer uth.lock.Unlock()
	if tokData, ok := uth.activeTokens[token]; ok {
		// the token is valid. and lets compare the username as well to make sure the token is genuine
		if tokData.uname != uname {
			return "", 0, false, errors.New("token has invalid username provided.")
		}
		validity := time.Until(tokData.expiry)
		if validity <= 0 {
			return "", 0, false, errors.New("token is expired!!")
		}
		return tokData.uname, validity, tokData.isAdmin, nil
	}
	return "", 0, false, errors.New("invalid token, cannot find it in the system")

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
		activeTokens:   make(map[string]*activeTokenData),
		envValues:      userEnvValues,
	}, nil
}
