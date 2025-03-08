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

// MockProducer is a mock implementation of storage.UserRepo
type MockProducer struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user storage.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*storage.User), args.Error(1)
}

func (m *MockUserRepo) GetUsers(ctx context.Context, limit rune) ([]storage.User, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]storage.User), args.Error(1)
}

func (m *MockProducer) PublishUser(user entity.User) error {
	args := m.Called(user)
	return args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockProducer := new(MockProducer)
	service := userservice.NewUserService(mockRepo, mockProducer)

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

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockProducer := new(MockProducer)
	service := userservice.NewUserService(mockRepo, mockProducer)

	mockRepo.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockProducer := new(MockProducer)
	service := userservice.NewUserService(mockRepo, mockProducer)

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
