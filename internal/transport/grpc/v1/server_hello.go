package grpcv1

import (
	"context"

	"github.com/pkg/errors"

	servicev1 "github.com/powerdigital/go-micro/internal/service/v1"
	pbgreeterv1 "github.com/powerdigital/go-micro/pkg/greeter/v1"
)

type ServerGreetingService interface {
	servicev1.GreetingService
}

type GreetingServer struct {
	service ServerGreetingService
	pbgreeterv1.UnimplementedGreeterServer
}

func NewGreetingServer(service servicev1.GreetingService) *GreetingServer {
	//nolint:exhaustruct
	return &GreetingServer{
		service: service,
	}
}

func (s *GreetingServer) GetHello(
	_ context.Context,
	req *pbgreeterv1.HelloRequest,
) (*pbgreeterv1.HelloResponse, error) {
	name := req.GetName()

	hello, err := s.service.GetHello(name)
	if err != nil {
		return nil, errors.Wrap(err, "get hello name")
	}

	return &pbgreeterv1.HelloResponse{
		Message: hello,
	}, nil
}
