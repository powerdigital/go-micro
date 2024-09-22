package service

import "github.com/pkg/errors"

type Service interface {
	GetHello(name string) (string, error)
}

type service struct{}

//nolint:revive
func NewService() *service {
	return &service{}
}

func (s *service) GetHello(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("the name provided to service is empty")
	}

	return "Hello, " + name, nil
}
