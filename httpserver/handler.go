package httpserver

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type handlerConfig struct {
	svc         service
	logger      log.Logger
	rateLimiter *rate.Limiter
}

// newHandler creates a new HTTP handler serving service endpoints.
func newHandler(cfg *handlerConfig) http.Handler {
	svc := &loggingMiddleware{next: cfg.svc, logger: cfg.logger}

	createUserEndpoint := makeCreateUserEndpoint(svc)
	createUserEndpoint = applyMiddleware(createUserEndpoint, "CreateUser", cfg)

	authUserEndpoint := makeAuthUserEndpoint(svc)
	authUserEndpoint = applyMiddleware(authUserEndpoint, "AuthUser", cfg)

	router := mux.NewRouter()

	router.Path("/api/v1/user").Methods("POST").Handler(kithttp.NewServer(
		createUserEndpoint,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	))

	router.Path("/api/v1/login/{login}/password/{password}").Methods("GET").Handler(kithttp.NewServer(
		authUserEndpoint,
		decodeAuthUserRequest,
		encodeAuthUserResponse,
	))

	return router
}

func applyMiddleware(e endpoint.Endpoint, method string, cfg *handlerConfig) endpoint.Endpoint {
	return ratelimit.NewErroringLimiter(cfg.rateLimiter)(e)
}
