package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"

	grpcv1 "github.com/powerdigital/go-micro/internal/transport/grpc/v1"
	hellov1 "github.com/powerdigital/go-micro/pkg/grpc/v1"
)

func (b *Builder) GRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	greetingServer, err := b.GreetingServer()
	if err != nil {
		return nil, errors.Wrap(err, "build greeting server")
	}

	hellov1.RegisterGreeterAPIServer(grpcServer, greetingServer)

	b.shutdown.add(func(_ context.Context) error {
		grpcServer.GracefulStop()

		return nil
	})

	return grpcServer, nil
}

func (b *Builder) GreetingServer() (*grpcv1.GRPCHandler, error) {
	service, err := b.GreetingService()
	if err != nil {
		return nil, errors.Wrap(err, "build GreetingService")
	}

	return grpcv1.NewGRPCHandler(service), nil
}
