package kafkav1

import (
	"context"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"
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

func (h *ConsumerGroupHandler) Healthcheck(_ context.Context) error {
	if h.running.Load() {
		return nil
	}

	return errors.New("handler is not running")
}
