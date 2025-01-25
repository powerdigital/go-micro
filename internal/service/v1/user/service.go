package userservice

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

type UserSrv interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (entity.User, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, userID int64) error
}

type UserService struct {
	repo storage.UserRepo
}

func NewUserService(repo storage.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	id, err := s.repo.CreateUser(ctx, user.EntityToModel())
	if err != nil {
		return 0, errors.Wrap(err, "create user")
	}

	return id, nil
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (entity.User, error) {
	userModel, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "user not found")
	}

	return entity.ModelToEntity(userModel), nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	userModels, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get users")
	}

	users := make([]entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = entity.ModelToEntity(&userModel)
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user entity.User) error {
	if err := s.repo.UpdateUser(ctx, user.EntityToModel()); err != nil {
		return errors.Wrap(err, "user service UpdateUser")
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return errors.Wrap(err, "user service DeleteUser")
	}

	return nil
}
