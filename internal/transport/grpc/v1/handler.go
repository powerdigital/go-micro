package grpcv1

import (
	"context"

	"github.com/cockroachdb/errors"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	userv1 "github.com/powerdigital/go-micro/pkg/grpc/v1"
)

type ServerUserService interface {
	userservice.UserSrv
}

type GRPCHandler struct {
	service ServerUserService
	userv1.UnimplementedUserAPIServer
}

func NewGRPCHandler(service userservice.UserSrv) *GRPCHandler {
	//nolint:exhaustruct
	return &GRPCHandler{
		service: service,
	}
}

func (s *GRPCHandler) CreateUser(
	ctx context.Context,
	req *userv1.CreateUserRequest,
) (*userv1.CreateUserResponse, error) {
	//nolint:exhaustruct
	user := entity.User{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Phone: req.GetPhone(),
		Age:   int(req.GetAge()),
	}

	userID, err := s.service.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "create user")
	}

	//nolint:gosec
	return &userv1.CreateUserResponse{
		UserId: uint32(userID),
	}, nil
}
