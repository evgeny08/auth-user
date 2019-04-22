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

// CreateSession create new session in storage
func (s *Storage) CreateSession(ctx context.Context, session *types.Session) error {
	err := db.C(collectionAuth).Insert(&session)
	return err
}

// FindAccessToken find AccessToken by storage
func (s *Storage) FindAccessToken(ctx context.Context, clientToken string) (*types.Session, error) {
	var session *types.Session
	filter := bson.M{"access_token": clientToken}
	err := db.C(collectionAuth).Find(filter).One(&session)
	return session, err
}
