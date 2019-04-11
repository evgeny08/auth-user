package httpserver

import (
	"context"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/evgeny08/auth-user/types"
)

type mockService struct {
	onCreateUser func(ctx context.Context, user *types.User) error
}

func (s *mockService) createUser(ctx context.Context, user *types.User) error {
	return s.onCreateUser(ctx, user)
}

func startTestServer(t *testing.T) (*httptest.Server, *Client, *mockService) {
	svc := &mockService{}

	handler := newHandler(&handlerConfig{
		svc:         svc,
		logger:      log.NewNopLogger(),
		rateLimiter: rate.NewLimiter(rate.Inf, 1),
	})

	server := httptest.NewServer(handler)

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	return server, client, svc
}

func TestCreateStructure(t *testing.T) {
	server, client, svc := startTestServer(t)
	defer server.Close()

	testCases := []struct {
		name string
		user *types.User
		err  error
	}{
		{
			name: "ok response",
			user: &types.User{
				Login:    "qww",
				Password: "12345678",
			},
			err: nil,
		},
		{
			name: "err response",
			user: &types.User{
				Login:    "",
				Password: "12345678",
			},
			err: errorf(ErrBadParams, "empty login"),
		},
		{
			name: "err response",
			user: &types.User{
				Login:    "111",
				Password: "11",
			},
			err: errorf(ErrBadParams, "password must be 8 symbols min"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onCreateUser = func(ctx context.Context, user *types.User) error {
				return tc.err
			}
			gotErr := client.CreateUser(context.Background(), tc.user)
			if !reflect.DeepEqual(gotErr, tc.err) {
				t.Fatalf("got error %#v want %#v", gotErr, tc.err)
			}

		})
	}
}
