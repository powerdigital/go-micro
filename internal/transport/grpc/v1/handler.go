package grpcv1

import (
	"context"

	"github.com/cockroachdb/errors"

	servicev1 "github.com/powerdigital/go-micro/internal/service/v1/greeting"
	grpcv1 "github.com/powerdigital/go-micro/pkg/grpc/v1"
)

type ServerGreetingService interface {
	servicev1.HelloService
}

type GRPCHandler struct {
	service ServerGreetingService
	grpcv1.UnimplementedGreeterAPIServer
}

func NewGRPCHandler(service servicev1.HelloService) *GRPCHandler {
	//nolint:exhaustruct
	return &GRPCHandler{
		service: service,
	}
}

func (s *GRPCHandler) GetHello(
	_ context.Context,
	req *grpcv1.GetHelloRequest,
) (*grpcv1.GetHelloResponse, error) {
	name := req.GetName()

	hello, err := s.service.GetHello(name)
	if err != nil {
		return nil, errors.Wrap(err, "get hello name")
	}

	return &grpcv1.GetHelloResponse{
		Message: hello,
	}, nil
}
