package userservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
)

// MockUserRepo is a mock implementation of storage.UserRepo
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*storage.User), args.Error(1)
}

func (m *MockUserRepo) GetUsers(ctx context.Context) ([]storage.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]storage.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user storage.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := userservice.NewUserService(mockRepo)

	user := entity.User{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}

	mockRepo.On("CreateUser", mock.Anything, user.EntityToModel()).Return(int64(1), nil)

	id, err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := userservice.NewUserService(mockRepo)

	User := &storage.User{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}
	userEntity := entity.ModelToEntity(User)

	mockRepo.On("GetUser", mock.Anything, int64(1)).Return(User, nil)

	user, err := service.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, userEntity, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := userservice.NewUserService(mockRepo)

	mockRepo.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
