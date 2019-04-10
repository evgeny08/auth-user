package storage

import (
	"context"

	"github.com/evgeny08/auth-user/types"
)

// CreateUser creates a user in storage
func (s *Storage) CreateUser(ctx context.Context, user *types.User) error {
	err := db.C(collectionUser).Insert(&user)
	return err
}
