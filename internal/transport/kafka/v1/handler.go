package kafkav1

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user"
	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type Handler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type ConsumerGroupHandler struct {
	Handler Handler
	running atomic.Bool
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.running.Store(true)

	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.running.Store(false)

	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	messages := claim.Messages()

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return nil
			}

			if err := h.Handler.Handle(ctx, message); err != nil {
				return errors.Wrap(err, "handle message")
			}

			session.MarkMessage(message, "")
		case <-ctx.Done():
			return nil
		}
	}
}

type UserCreatingConsumer struct {
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
