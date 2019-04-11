package storage

import (
	"context"
	"errors"
	"sync"

	"gopkg.in/mgo.v2"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var db *mgo.Database

const (
	collectionUser = "auth_user"
	collectionAuth = "auth_main"
)

// Storage stores keys.
type Storage struct {
	url    string
	dbName string
	logger log.Logger

	mu      sync.RWMutex
	session *mgo.Session // Master session.
	lastErr error

	ctx    context.Context
	cancel context.CancelFunc
	donec  chan struct{}
}

// Config is a storage configuration.
type Config struct {
	URL    string
	Logger log.Logger
	DBName string
}

// New creates a new MongoDB storage using the given configuration.
func New(cfg *Config) (*Storage, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Storage{
		url:    cfg.URL,
		dbName: cfg.DBName,
		logger: cfg.Logger,

		ctx:    ctx,
		cancel: cancel,
		donec:  make(chan struct{}),
	}

	err := s.connect(cfg)
	if err != nil {
		return nil, level.Error(s.logger).Log("msg", "failed to connect mongodb", "error:", err)
	}
	return s, nil
}

func (s *Storage) connect(cfg *Config) error {
	defer close(s.donec)
	for {
		// Check if we're canceled.
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		session, err := mgo.Dial(cfg.URL)
		if err != nil {
			return err
		}

		db = session.DB(cfg.DBName)

		if err != nil {
			// Check if we're canceled
			// once more before sleeping.
			select {
			case <-s.ctx.Done():
				return nil
			default:
			}
			s.logger.Log("failed to connect to mongo: %v", err)
			continue
		}
		s.logger.Log("msg", "established mongo connection")
		s.mu.Lock()
		s.session = session
		s.mu.Unlock()
		return nil
	}
}

// Shutdown close mongo session
func (s *Storage) Shutdown() {
	// Cancel connect loop.
	s.cancel()
	<-s.donec

	// Close mongo session.
	s.mu.Lock()
	if s.session != nil {
		s.session.Close()
		s.session = nil
		s.lastErr = errors.New("mongoclient is shut down")
	}
	s.mu.Unlock()

	level.Info(s.logger).Log("msg", "mongoclient: shutdown complete")
}
