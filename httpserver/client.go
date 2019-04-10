package httpserver

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/evgeny08/auth-user/types"
)

// Client is a client for Key Generator service.
type Client struct {
	createUser endpoint.Endpoint
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
