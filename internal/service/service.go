package service

type Service interface {
	GetHello(name string) (string, error)
}

type service struct{}

//nolint:revive
func NewService() *service {
	return &service{}
}

func (s *service) GetHello(name string) (string, error) {
	return "Hello, " + name, nil
}
