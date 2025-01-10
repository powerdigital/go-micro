package greeting

import (
	"github.com/cockroachdb/errors"
)

type HelloSrv interface {
	GetHello(name string) (string, error)
}

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) GetHello(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("the name provided to the HelloService is empty")
	}

	return "Hello, " + name, nil
}
