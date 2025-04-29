package producer

import (
	"encoding/json"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"

	"github.com/powerdigital/go-micro/internal/service/v1/user/entity"
)

type UserQueue interface {
	PublishUser(user entity.User) error
}

type Producer struct {
	SyncProducer    sarama.SyncProducer
	CreateUserTopic string
}

func (p *Producer) PublishUser(user entity.User) error {
	userMsg, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "marshal user message")
	}

	//nolint:exhaustruct
	message := &sarama.ProducerMessage{
		Topic: p.CreateUserTopic,
		Key:   sarama.StringEncoder(strconv.FormatInt(user.ID, 10)),
		Value: sarama.StringEncoder(userMsg),
	}

	_, _, err = p.SyncProducer.SendMessage(message)
	if err != nil {
		return errors.Wrap(err, "send user message")
	}

	return nil
}
