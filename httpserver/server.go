package httpserver

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"

	"github.com/evgeny08/auth-user/types"
)

// ServerHTTP is a service structure http server.
type ServerHTTP struct {
	logger log.Logger
	srv    *http.Server
}

// Config is a http server configuration.
type Config struct {
	Logger      log.Logger
	Port        string
	Storage     Storage
	RateLimiter *rate.Limiter
	ServerNATS  ServerNATS
	WebSocket   WebSocket
}

// Storage is a persistent auth-user storage.
type Storage interface {
	CreateUser(ctx context.Context, user *types.User) error
	FindUserByLogin(ctx context.Context, login string) (*types.User, error)
	CreateSession(ctx context.Context, session *types.Session) error
	FindAccessToken(ctx context.Context, clientToken string) (*types.Session, error)
}

type ServerNATS interface {
	Send(msg string) error
}

type WebSocket interface {
	WsHandler(w http.ResponseWriter, r *http.Request)
}

// New creates a new http server.
func New(cfg *Config) (*ServerHTTP, error) {
	mu := http.NewServeMux()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mu,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := &ServerHTTP{
		logger: cfg.Logger,
		srv:    srv,
	}

	svc := &basicService{
		logger:     cfg.Logger,
		storage:    cfg.Storage,
		serverNATS: cfg.ServerNATS,
		webSocket:  cfg.WebSocket,
	}

	handler := newHandler(&handlerConfig{
		svc:         svc,
		logger:      cfg.Logger,
		rateLimiter: cfg.RateLimiter,
	})

	mu.Handle("/api/v1/", accessControl(handler))

	router := mux.NewRouter()
	router.HandleFunc("/ws", svc.webSocket.WsHandler)

	http.ListenAndServe(":8844", router)

	return server, nil
}

// CORS headers
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Run starts the server.
func (s *ServerHTTP) Run() error {
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown stopped the http server.
func (s *ServerHTTP) Shutdown() {
	err := s.srv.Close()
	if err != nil {
		err := level.Info(s.logger).Log("msg", "HTTP server: shutdown has err", "err:", err)
		if err != nil {
			return
		}
	}
	err = level.Info(s.logger).Log("msg", "HTTP server: shutdown complete")
	if err != nil {
		return
	}
}
