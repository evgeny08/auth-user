package httpserver

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/evgeny08/auth-user/types"
)

// Client is a client for auth-user service.
type Client struct {
	createUser      endpoint.Endpoint
	authUser        endpoint.Endpoint
	findUserByLogin endpoint.Endpoint
}

// NewClient creates a new service client.
func NewClient(serviceURL string) (*Client, error) {
	baseURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		createUser: kithttp.NewClient(
			"POST",
			baseURL,
			encodeCreateUserRequest,
			decodeCreateUserResponse,
		).Endpoint(),

		authUser: kithttp.NewClient(
			"GET",
			baseURL,
			encodeAuthUserRequest,
			decodeAuthUserResponse,
		).Endpoint(),

		findUserByLogin: kithttp.NewClient(
			"GET",
			baseURL,
			encodeFindUserByLoginRequest,
			decodeFindUserByLoginResponse,
		).Endpoint(),
	}

	return c, nil
}

// CreateUser creates a new user.
func (c *Client) CreateUser(ctx context.Context, user *types.User) error {
	request := createUserRequest{User: user}
	response, err := c.createUser(ctx, request)
	if err != nil {
		return err
	}
	res := response.(createUserResponse)
	return res.Err
}

// AuthUser Authentication user into the system.
func (c *Client) AuthUser(ctx context.Context, login, password string) (*types.Session, error) {
	request := authUserRequest{Login: login, Password: password}
	response, err := c.authUser(ctx, request)
	if err != nil {
		return nil, err
	}
	res := response.(authUserResponse)
	return res.Session, res.Err
}

// findUserByLogin return user by login
func (c *Client) FindUserByLogin(ctx context.Context, login string) (*types.User, error) {
	request := findUserByLoginRequest{Login: login}
	response, err := c.findUserByLogin(ctx, request)
	if err != nil {
		return nil, err
	}
	res := response.(findUserByLoginResponse)
	return res.User, res.Err
}
