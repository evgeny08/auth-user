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
