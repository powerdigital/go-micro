package build

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

type Consumer struct {
	Handler       sarama.ConsumerGroupHandler
	ConsumerGroup sarama.ConsumerGroup
	Topic         string
}

func (c *Consumer) Consume(ctx context.Context) error {
	return c.ConsumerGroup.Consume(ctx, []string{c.Topic}, c.Handler) //nolint:wrapcheck
}

func (b *Builder) consumerGroup(ctx context.Context, group string) (sarama.ConsumerGroup, error) { //nolint:ireturn
	sarama.Logger = zerolog.Ctx(ctx)

	cfg := sarama.NewConfig()

	cfg.ClientID = b.config.App.Name
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.MaxProcessingTime = time.Minute * 5 //nolint:mnd
	cfg.Consumer.Return.Errors = true

	cg, err := sarama.NewConsumerGroup(b.config.Kafka.Brokers, group, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create consumer group")
	}

	go func() {
		for err := range cg.Errors() {
			zerolog.Ctx(ctx).Err(err).Msgf("%s consumer group", group)
		}
	}()

	return cg, nil
}
