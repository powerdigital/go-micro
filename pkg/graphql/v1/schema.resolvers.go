//nolint:gosec
package v1

import (
	"context"
	"fmt"

	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, user CreateUser) (int32, error) {
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

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, user UpdateUser) (int32, error) {
	userEntity := entity.User{
		ID:    int64(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   int(user.Age),
	}

	err := r.service.UpdateUser(ctx, userEntity)
	if err != nil {
		return 0, fmt.Errorf("failed to update user: %w", err)
	}

	return user.ID, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id int32) (int32, error) {
	err := r.service.DeleteUser(ctx, int64(id))
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}

	return id, nil
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id int32) (*User, error) {
	u, err := r.service.GetUser(ctx, int64(id))
	if err != nil {
		return &User{}, fmt.Errorf("failed to get user: %w", err)
	}

	user := User{
		ID:    int32(u.ID),
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
		Age:   int32(u.Age),
	}

	return &user, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context) ([]*User, error) {
	users, err := r.service.GetUsers(ctx)
	if err != nil {
		return []*User{}, fmt.Errorf("failed to get users: %w", err)
	}

	userList := make([]*User, 0, len(users))

	for _, u := range users {
		user := User{
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
