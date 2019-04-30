package natsserver

import (
	"github.com/nats-io/go-nats"
	"log"
	"sync"
)

type ServerNATS struct {
	logger log.Logger
	srv    *nats.Conn
}

type Config struct {
	Logger log.Logger
	URL    string
}

func New(cfg *Config) (*ServerNATS, error) {
	// Create NATS server connection
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	srv := &ServerNATS{
		logger: log.Logger{},
		srv:    nc,
	}
	return srv, nil
}

func (s *ServerNATS) Run() error {
	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	if _, err := s.srv.Subscribe("updates", func(m *nats.Msg) {
		log.Printf("%s: %s", m.Subject, m.Data)
	}); err != nil {
		log.Fatal(err)
	}
	// Wait for a message to come in
	wg.Wait()
	return nil
}
