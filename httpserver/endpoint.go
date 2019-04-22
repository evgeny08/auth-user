package httpserver

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"github.com/evgeny08/auth-user/types"
)

func makeCreateUserEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		err := svc.createUser(ctx, req.User)
		return createUserResponse{Err: err}, nil
	}
}

type createUserRequest struct {
	User *types.User
}

type createUserResponse struct {
	Err error
}

func makeAuthUserEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authUserRequest)
		session, err := svc.authUser(ctx, req.Login, req.Password)
		return authUserResponse{Session: session, Err: err}, nil
	}
}

type authUserRequest struct {
	Login    string
	Password string
}

type authUserResponse struct {
	Session *types.Session
	Err     error
}

func makeFindUserByLoginEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findUserByLoginRequest)
		user, err := svc.findUserByLogin(ctx, req.Login, req.ClientToken)
		return findUserByLoginResponse{User: user, Err: err}, nil
	}
}

type findUserByLoginRequest struct {
	Login       string
	ClientToken string `json:"client_token"`
}

type findUserByLoginResponse struct {
	User *types.User
	Err  error
}
