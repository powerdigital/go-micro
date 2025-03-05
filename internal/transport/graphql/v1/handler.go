//nolint:gosec
package graphqlv1

import (
	"context"
	"fmt"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

type GqlHandler struct {
	service userservice.UserSrv
}

func NewGqlHandler(service userservice.UserSrv) *GqlHandler {
	return &GqlHandler{
		service: service,
	}
}

func (r *GqlHandler) CreateUser(ctx context.Context, user UserRequest) (int32, error) {
	//nolint:exhaustruct
	userEntity := entity.User{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   int(user.Age),
	}

	userID, err := r.service.CreateUser(ctx, userEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return int32(userID), nil
}

func (r *GqlHandler) UpdateUser(ctx context.Context, user UserRequest) (int32, error) {
	userID := *user.ID
	userEntity := entity.User{
		ID:    int64(userID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   int(user.Age),
	}

	err := r.service.UpdateUser(ctx, userEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update user: %w", err)
	}

	return userID, nil
}

func (r *GqlHandler) DeleteUser(ctx context.Context, id int32) (int32, error) {
	err := r.service.DeleteUser(ctx, int64(id))
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}

	return id, nil
}

func (r *GqlHandler) GetUser(ctx context.Context, id int32) (*UserResponse, error) {
	u, err := r.service.GetUser(ctx, int64(id))
	if err != nil {
		return &UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	user := UserResponse{
		ID:    int32(u.ID),
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   int32(u.Age),
	}

	return &user, nil
}

func (r *GqlHandler) GetUsers(ctx context.Context, limit int32) ([]*UserResponse, error) {
	users, err := r.service.GetUsers(ctx, limit)
	if err != nil {
		return []*UserResponse{}, fmt.Errorf("failed to get users: %w", err)
	}

	userList := make([]*UserResponse, 0, len(users))

	for _, u := range users {
		user := UserResponse{
			ID:    int32(u.ID),
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
			Age:   int32(u.Age),
		}

		userList = append(userList, &user)
	}

	return userList, nil
}
