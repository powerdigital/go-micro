package build

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

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
