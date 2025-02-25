package build

import (
	"context"

	kafkav1 "github.com/powerdigital/go-micro/internal/transport/kafka/v1"
)

func (b *Builder) UserCreatingConsumer(ctx context.Context) (*Consumer, error) {
	handler, err := b.userCreatingHandler(ctx)
	if err != nil {
		return nil, err
	}

	topic := b.config.Kafka.TopicCreateUser

	group, err := b.consumerGroup(ctx, b.config.Kafka.ConsumerGroup)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Handler:       handler,
		ConsumerGroup: group,
		Topic:         topic,
	}, nil
}

func (b *Builder) userCreatingHandler(ctx context.Context) (*kafkav1.ConsumerGroupHandler, error) {
	service, err := b.UserService(ctx)
	if err != nil {
		return nil, err
	}

	sub := kafkav1.UserCreatingConsumer{
		UserService: service,
	}

	handler := kafkav1.ConsumerGroupHandler{
		Handler: &sub,
	}

	return &handler, nil
}
