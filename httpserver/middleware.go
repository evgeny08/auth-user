package httpserver

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/evgeny08/auth-user/types"
)

// loggingMiddleware wraps Service and logs request information to the provided logger.
type loggingMiddleware struct {
	next   service
	logger log.Logger
}

func (m *loggingMiddleware) createUser(ctx context.Context, user *types.User) error {
	begin := time.Now()
	err := m.next.createUser(ctx, user)
	level.Info(m.logger).Log(
		"method", "CreateUser",
		"err", err,
		"elapsed", time.Since(begin),
	)
	return err
}

func (m *loggingMiddleware) authUser(ctx context.Context, login, password string) (*types.Session, error) {
	begin := time.Now()
	session, err := m.next.authUser(ctx, login, password)
	level.Info(m.logger).Log(
		"method", "AuthUser",
		"err", err,
		"elapsed", time.Since(begin),
	)
	return session, err
}

func (m *loggingMiddleware) findUserByLogin(ctx context.Context, login string) (*types.User, error) {
	begin := time.Now()
	user, err := m.next.findUserByLogin(ctx, login)
	level.Info(m.logger).Log(
		"method", "FindUserByLogin",
		"err", err,
		"elapsed", time.Since(begin),
	)
	return user, err
}
