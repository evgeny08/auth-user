package httpserver

import (
	"context"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/evgeny08/auth-user/types"
)

// service manages HTTP server methods.
type service interface {
	createUser(ctx context.Context, user *types.User) error
	authUser(ctx context.Context, login, password string) (*types.Session, error)
	findUserByLogin(ctx context.Context, login, clientToken string) (*types.User, error)
}

type basicService struct {
	logger     log.Logger
	storage    Storage
	serverNATS ServerNATS
}

// createUser creates a new User
func (s *basicService) createUser(ctx context.Context, user *types.User) error {
	// Validate
	passwordLength := 8
	if strings.TrimSpace(user.Login) == "" {
		return errorf(ErrBadParams, "empty login")
	}
	if len(user.Password) < passwordLength {
		return errorf(ErrBadParams, "password must be 8 symbols min")
	}

	_, err := s.storage.FindUserByLogin(ctx, user.Login)
	if err == nil {
		return errorf(ErrConflict, "user with this login already exist")
	}

	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return errorf(ErrInternal, "failed to insert user: %v", err)
	}

	msg := "user created"
	err = s.serverNATS.Send(msg)
	if err != nil {
		return errorf(ErrInternal, "failed to send msg", err)
	}
	return nil
}

// authUser Authentication user into the system
func (s *basicService) authUser(ctx context.Context, login, password string) (*types.Session, error) {
	// Validate
	u, err := s.storage.FindUserByLogin(ctx, login)
	if err != nil {
		if storageErrIsNotFound(err) {
			return nil, errorf(ErrNotFound, "user with this login not found")
		}
		return nil, errorf(ErrInternal, "find user error: %v", err)
	}
	if u.Password != password {
		return nil, errorf(ErrBadParams, "wrong password")
	}
	session := &types.Session{
		AccessToken:         uuid.New().String(),
		ExpiresAccessToken:  time.Now().UTC().Add(30 * time.Minute).UnixNano(),
		RefreshToken:        uuid.New().String(),
		ExpiresRefreshToken: time.Now().UTC().AddDate(0, 2, 0).UnixNano(),
	}
	err = s.storage.CreateSession(ctx, session)
	if err != nil {
		return nil, errorf(ErrInternal, "failed to create session: %v", err)
	}
	return session, nil
}

// findUserByLogin find user in storage by login
func (s *basicService) findUserByLogin(ctx context.Context, login, clientToken string) (*types.User, error) {
	// Validate access Token
	sess, err := s.storage.FindAccessToken(ctx, clientToken)
	if err != nil {
		if storageErrIsNotFound(err) {
			return nil, errorf(ErrNotFound, "token not found")
		}
		return nil, errorf(ErrInternal, "find token error: %v", err)
	}
	if sess.AccessToken != clientToken {
		return nil, errorf(ErrInternal, "failed to authorisation: %v", err)
	}
	if sess.ExpiresAccessToken < time.Now().UTC().UnixNano() {
		return nil, errorf(ErrInternal, "the token has expired, log in")
	}
	// Validate
	if strings.TrimSpace(login) == "" {
		return nil, errorf(ErrBadParams, "empty login")
	}
	user, err := s.storage.FindUserByLogin(ctx, login)
	if err != nil {
		if storageErrIsNotFound(err) {
			return nil, errorf(ErrNotFound, "user is not found")
		}
		return nil, errorf(ErrInternal, "failed to find user: %v", err)
	}
	return user, nil
}

// storageErrIsNotFound checks if the storage error is "not found".
func storageErrIsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}
	e, ok := err.(notFound)
	return ok && e.NotFound()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
