//nolint:depguard
package build

import (
	"context"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"

	grpcv1 "github.com/powerdigital/go-micro/internal/transport/grpc/v1"
	userv1 "github.com/powerdigital/go-micro/pkg/grpc/v1"
)

func (b *Builder) GRPCServer(ctx context.Context) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	service, err := b.UserService(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "build UserService")
	}

	userServer := grpcv1.NewGRPCHandler(service)

	userv1.RegisterUserAPIServer(grpcServer, userServer)

	b.shutdown.add(func(_ context.Context) error {
		grpcServer.GracefulStop()

		return nil
	})

	return grpcServer, nil
}
