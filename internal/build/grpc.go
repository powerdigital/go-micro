package build

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	grpcv1 "github.com/powerdigital/go-micro/internal/transport/grpc/v1"
	greeterv1 "github.com/powerdigital/go-micro/pkg/greeter/v1"
)

func (b *Builder) GRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	greetingServer, err := b.GreetingServer()
	if err != nil {
		return nil, errors.Wrap(err, "build greeting server")
	}

	greeterv1.RegisterGreeterServer(grpcServer, greetingServer)

	b.shutdown.add(func(_ context.Context) error {
		grpcServer.GracefulStop()

		return nil
	})

	return grpcServer, nil
}

func (b *Builder) GreetingServer() (*grpcv1.GRPCHandler, error) {
	service, err := b.GreetingService()
	if err != nil {
		return nil, errors.Wrap(err, "build service")
	}

	return grpcv1.NewGRPCHandler(service), nil
}
