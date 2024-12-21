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

type GRPCHandler struct {
	service ServerGreetingService
	pbgreeterv1.UnimplementedGreeterAPIServer
}

func NewGRPCHandler(service servicev1.GreetingService) *GRPCHandler {
	//nolint:exhaustruct
	return &GRPCHandler{
		service: service,
	}
}

func (s *GRPCHandler) GetHello(
	_ context.Context,
	req *pbgreeterv1.GetHelloRequest,
) (*pbgreeterv1.GetHelloResponse, error) {
	name := req.GetName()

	hello, err := s.service.GetHello(name)
	if err != nil {
		return nil, errors.Wrap(err, "get hello name")
	}

	return &pbgreeterv1.GetHelloResponse{
		Message: hello,
	}, nil
}
