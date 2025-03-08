package config

type Kafka struct {
	Brokers          []string `envconfig:"KAFKA_BROKERS"`
	User             string   `envconfig:"KAFKA_USER"`
	Password         string   `envconfig:"KAFKA_PASSWORD"`
	ConsumerGroup    string   `envconfig:"KAFKA_CONSUMER_GROUP"`
	TopicCreateUser  string   `envconfig:"KAFKA_TOPIC_CREATE_USER"`
	TopicUserCreated string   `envconfig:"KAFKA_TOPIC_USER_CREATED"`
	TopicDeleteUser  string   `envconfig:"KAFKA_TOPIC_DELETE_USER"`
}
