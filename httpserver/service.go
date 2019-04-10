package httpserver

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/evgeny08/auth-user/types"
)

// service manages HTTP server methods.
type service interface {
	createUser(ctx context.Context, user *types.User) error
}

type basicService struct {
	logger  log.Logger
	storage Storage
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

	err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return errorf(ErrInternal, "failed to insert user: %v", err)
	}
	return nil
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

