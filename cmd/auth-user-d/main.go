package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/time/rate"

	"github.com/evgeny08/auth-user/httpserver"
	"github.com/evgeny08/auth-user/storage"
)

type configuration struct {
	HTTPPort       string        `envconfig:"AUTH_HTTP_PORT" default:"24020"`
	RateLimitEvery time.Duration `envconfig:"AUTH_RATE_LIMIT_EVERY" default:"1us"`
	RateLimitBurst int           `envconfig:"AUTH_RATE_LIMIT_BURST" default:"100"`

	MongoURL string `envconfig:"AUTH_MONGO_URL" default:"127.0.0.1:27017"`
	DBName   string `envconfig:"AUTH_DB_NAME"   default:"auth-user"`
}

func main() {
	const (
		exitCodeSuccess = 0
		exitCodeFailure = 1
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(exitCodeFailure)
	}

	mongoDB, err := storage.New(&storage.Config{
		URL:    cfg.MongoURL,
		DBName: cfg.DBName,
		Logger: logger,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize MongoDB", "err", err)
		os.Exit(exitCodeFailure)
	}

	serverHTTP, err := httpserver.New(&httpserver.Config{
		Logger:      logger,
		Port:        cfg.HTTPPort,
		Storage:     mongoDB,
		RateLimiter: rate.NewLimiter(rate.Every(cfg.RateLimitEvery), cfg.RateLimitBurst),
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize HTTP server", "err", err)
		os.Exit(exitCodeFailure)
	}
	go func() {
		level.Info(logger).Log("msg", "starting HTTP server", "port", cfg.HTTPPort)
		if err := serverHTTP.Run(); err != nil {
			level.Error(logger).Log("msg", "HTTP server run failure", "err", err)
			os.Exit(exitCodeFailure)
		}
	}()

	errc := make(chan error, 1)
	donec := make(chan struct{})
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	go func() {
		select {
		case sig := <-sigc:
			level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)
			serverHTTP.Shutdown() // Shutdown server HTTP
			mongoDB.Shutdown()    // Shutdown MongoDB
			signal.Stop(sigc)
			close(donec)
		case <-errc:
			level.Info(logger).Log("msg", "now exiting with error", "error code", exitCodeFailure)
			os.Exit(exitCodeFailure)
		}
	}()

	<-donec
	level.Info(logger).Log("msg", "goodbye")
	os.Exit(exitCodeSuccess)

}
