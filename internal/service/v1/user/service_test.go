package userservice

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql/mock"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage/mysql/model"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepo(ctrl)
	service := NewUserService(mockRepo)

	user := entity.User{Name: "John", Email: "john@example.com", Phone: "1234567890", Age: 30}
	modelUser := user.EntityToModel()

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), modelUser).
		Return(int64(1), nil)

	id, err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepo(ctrl)
	service := NewUserService(mockRepo)

	modelUser := model.User{ID: 1, Name: "John", Email: "john@example.com", Phone: "1234567890", Age: 30}
	expectedUser := entity.User{ID: 1, Name: "John", Email: "john@example.com", Phone: "1234567890", Age: 30}

	mockRepo.EXPECT().
		GetUser(gomock.Any(), int64(1)).
		Return(&modelUser, nil)

	user, err := service.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepo(ctrl)
	service := NewUserService(mockRepo)

	modelUsers := []model.User{
		{ID: 1, Name: "John", Email: "john@example.com", Phone: "1234567890", Age: 30},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Phone: "0987654321", Age: 25},
	}
	expectedUsers := []entity.User{
		{ID: 1, Name: "John", Email: "john@example.com", Phone: "1234567890", Age: 30},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Phone: "0987654321", Age: 25},
	}

	mockRepo.EXPECT().
		GetUsers(gomock.Any()).
		Return(modelUsers, nil)

	users, err := service.GetUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepo(ctrl)
	service := NewUserService(mockRepo)

	user := entity.User{ID: 1, Name: "Updated", Email: "updated@example.com", Phone: "9999999999", Age: 35}
	modelUser := user.EntityToModel()

	mockRepo.EXPECT().
		UpdateUser(gomock.Any(), modelUser).
		Return(nil)

	err := service.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepo(ctrl)
	service := NewUserService(mockRepo)

	mockRepo.EXPECT().
		DeleteUser(gomock.Any(), int64(1)).
		Return(nil)

	err := service.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
}
