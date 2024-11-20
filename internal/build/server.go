package build

import (
	"github.com/pkg/errors"

	grpcv1 "github.com/powerdigital/go-micro/internal/transport/grpc/v1"
)

func (b *Builder) GreetingServer() (*grpcv1.GreetingServer, error) {
	service, err := b.GreetingService()
	if err != nil {
		return nil, errors.Wrap(err, "build service")
	}

	return grpcv1.NewGreetingServer(service), nil
}
