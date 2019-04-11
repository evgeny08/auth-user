package storage

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/evgeny08/auth-user/types"
)

// CreateUser creates a user in storage
func (s *Storage) CreateUser(ctx context.Context, user *types.User) error {
	err := db.C(collectionUser).Insert(&user)
	return err
}

// FindUserByLogin find user by given login
func (s *Storage) FindUserByLogin(ctx context.Context, login string) (*types.User, error) {
	var user *types.User
	filter := bson.M{"login": login}
	err := db.C(collectionUser).Find(filter).One(&user)
	return user, err
}
