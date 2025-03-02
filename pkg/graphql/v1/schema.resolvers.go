//nolint:gosec
package v1

import (
	"context"
	"fmt"

	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
