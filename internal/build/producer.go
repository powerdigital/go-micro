package build

import (
	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"

	userservice "github.com/powerdigital/go-micro/internal/service/v1/user/producer"
)

func (b *Builder) Producer(brokers []string, topic string) (*userservice.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "producer creating")
	}

	return &userservice.Producer{
		SyncProducer:    producer,
		CreateUserTopic: topic,
	}, nil
}
