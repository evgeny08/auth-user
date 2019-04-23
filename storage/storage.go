package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"go.mongodb.org/mongo-driver/mongo"
)

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
	session *mongo.Database
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
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		session, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URL))
		if err != nil {
			return err
		}



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
		s.session = session.Database(cfg.DBName)
		s.mu.Unlock()
		return nil
	}
}

// Shutdown close mongo session
func (s *Storage) Shutdown() {
	// Cancel connect loop.
	s.cancel()
	<-s.donec

	err = s.session.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
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
