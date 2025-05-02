package userservice_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
	"github.com/powerdigital/go-micro/internal/service/v1/user/storage"
	storagemocks "github.com/powerdigital/go-micro/internal/service/v1/user/storage/mocks"
	producermocks "github.com/powerdigital/go-micro/pkg/producer/mocks"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx   context.Context
	repo  *storagemocks.UserRepo
	queue *producermocks.UserQueue
	srv   *userservice.UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.repo = &storagemocks.UserRepo{}
	suite.queue = &producermocks.UserQueue{}
	suite.srv = userservice.NewUserService(suite.repo, suite.queue)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.repo.AssertExpectations(suite.T())
	suite.queue.AssertExpectations(suite.T())
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreateUserSuccess() {
	user := entity.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "+1234567890",
		Age:   30,
	}

	suite.repo.On("CreateUser", suite.ctx, mock.Anything).Return(int64(123), nil)
	suite.queue.On("PublishUser", mock.Anything).Return(nil)

	id, err := suite.srv.CreateUser(suite.ctx, user)

	suite.NoError(err)
	suite.Equal(int64(123), id)

	expected, err := os.ReadFile("testdata/golden/create_user_success.golden.json")
	suite.Require().NoError(err)

	actual, err := json.Marshal(user)
	suite.Require().NoError(err)

	suite.JSONEq(string(expected), string(actual))
}

func (suite *UserServiceTestSuite) TestCreateUser_QueueFails() {
	user := entity.User{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Phone: gofakeit.Phone(),
		Age:   gofakeit.Number(18, 60),
	}

	suite.repo.On("CreateUser", suite.ctx, mock.Anything).Return(int64(321), nil)
	suite.queue.On("PublishUser", mock.Anything).Return(errors.New("kafka down"))

	id, err := suite.srv.CreateUser(suite.ctx, user)

	suite.Error(err)
	suite.Equal(int64(321), id)
}

func (suite *UserServiceTestSuite) TestUpdateUserSuccess() {
	user := entity.User{
		ID:    1,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Phone: gofakeit.Phone(),
		Age:   gofakeit.Number(18, 60),
	}

	suite.repo.On("UpdateUser", suite.ctx, mock.Anything).Return(nil)

	err := suite.srv.UpdateUser(suite.ctx, user)
	suite.NoError(err)
}

func (suite *UserServiceTestSuite) TestDeleteUserSuccess() {
	suite.repo.On("DeleteUser", suite.ctx, int64(42)).Return(nil)

	err := suite.srv.DeleteUser(suite.ctx, 42)
	suite.NoError(err)
}

func (suite *UserServiceTestSuite) TestGetUserNotFound() {
	suite.repo.On("GetUser", suite.ctx, int64(99)).Return(nil, errors.New("not found"))

	_, err := suite.srv.GetUser(suite.ctx, 99)
	suite.Error(err)
}

func (suite *UserServiceTestSuite) TestGetUsersEmpty() {
	suite.repo.On("GetUsers", suite.ctx, rune(0)).Return([]storage.User{}, nil)

	users, err := suite.srv.GetUsers(suite.ctx, rune(0))
	suite.NoError(err)
	suite.Len(users, 0)
}

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
