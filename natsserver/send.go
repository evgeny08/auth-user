package natsserver

import "log"

func (s *ServerNATS) Send(msg string) error {
	if err := s.srv.Publish("updates", []byte(msg)); err != nil {
		log.Fatal(err)
	}
	return nil
}
