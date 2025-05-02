package userservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
	storagemocks "github.com/powerdigital/go-micro/internal/service/v1/user/storage/mocks"
	producermocks "github.com/powerdigital/go-micro/pkg/producer/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	user := entity.User{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}

	mockRepo := new(storagemocks.UserRepo)
	mockRepo.On("CreateUser", mock.Anything, user.EntityToModel()).Return(int64(1), nil)

	mockProducer := new(producermocks.UserQueue)
	mockProducer.On("PublishUser", user).Return(nil)

	service := userservice.NewUserService(mockRepo, mockProducer)

	id, err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(storagemocks.UserRepo)
	mockProducer := new(producermocks.UserQueue)
	service := userservice.NewUserService(mockRepo, mockProducer)

	mockRepo.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(storagemocks.UserRepo)
	mockProducer := new(producermocks.UserQueue)
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
