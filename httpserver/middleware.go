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
	_ = level.Info(m.logger).Log(
		"method", "CreateUser",
		"err", err,
		"elapsed", time.Since(begin),
	)
	return err
}
