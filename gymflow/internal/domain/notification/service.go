package notification

import "log"

type Service interface {
	SendEmail(to, subject, body string) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) SendEmail(to, subject, body string) error {
	log.Printf("[EMAIL] to=%s subject=%s body=%s", to, subject, body)
	return nil
}
