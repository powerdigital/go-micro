package kafkav1

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

type UserCreatingConsumer struct {
	UserService userservice.UserSrv
}

type UserDeletingConsumer struct {
	UserService userservice.UserSrv
}

func (c *UserCreatingConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var userMsg entity.User

	err := json.Unmarshal(msg.Value, &userMsg)
	if err != nil {
		return errors.Wrap(err, "unmarshall user message")
	}

	_, err = c.UserService.CreateUser(ctx, userMsg)

	return errors.Wrap(err, "handle user message")
}

func (c *UserDeletingConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var userMsg entity.User

	err := json.Unmarshal(msg.Value, &userMsg)
	if err != nil {
		return errors.Wrap(err, "unmarshall user message")
	}

	err = c.UserService.DeleteUser(ctx, userMsg.ID)

	return errors.Wrap(err, "handle user message")
}
