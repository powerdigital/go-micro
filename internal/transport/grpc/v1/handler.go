//nolint:gosec
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

	return &userv1.CreateUserResponse{
		UserId: uint32(userID),
	}, nil
}

func (s *GRPCHandler) UpdateUser(
	ctx context.Context,
	req *userv1.UpdateUserRequest,
) (*userv1.UpdateUserResponse, error) {
	user := entity.User{
		ID:    int64(req.GetId()),
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Phone: req.GetPhone(),
		Age:   int(req.GetAge()),
	}

	err := s.service.UpdateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "update user")
	}

	return &userv1.UpdateUserResponse{
		UserId: req.GetId(),
	}, nil
}

func (s *GRPCHandler) DeleteUser(
	ctx context.Context,
	req *userv1.DeleteUserRequest,
) (*userv1.DeleteUserResponse, error) {
	userID := int64(req.GetUserId())

	err := s.service.DeleteUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "delete user")
	}

	return &userv1.DeleteUserResponse{
		UserId: req.GetUserId(),
	}, nil
}

func (s *GRPCHandler) GetUser(
	ctx context.Context,
	req *userv1.GetUserRequest,
) (*userv1.GetUserResponse, error) {
	userID := int64(req.GetUserId())

	user, err := s.service.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	return &userv1.GetUserResponse{
		Id:    uint32(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   uint32(user.Age),
	}, nil
}

func (s *GRPCHandler) GetUsers(
	ctx context.Context,
	req *userv1.GetUsersRequest,
) (*userv1.GetUsersResponse, error) {
	limit := rune(req.GetLimit())

	usersList, err := s.service.GetUsers(ctx, limit)
	if err != nil {
		return nil, errors.Wrap(err, "get users")
	}

	users := make([]*userv1.GetUserResponse, 0, len(usersList))

	for _, u := range usersList {
		user := userv1.GetUserResponse{
			Id:    uint32(u.ID),
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
			Age:   uint32(u.Age),
		}

		users = append(users, &user)
	}

	return &userv1.GetUsersResponse{
		Users: users,
	}, nil
}
