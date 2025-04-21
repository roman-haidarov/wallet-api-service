package kafka

import (
	// "encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"wallet-api-service/internal/config"

	// "wallet-api-service/internal/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/google/uuid"
	// "github.com/rs/zerolog/log"
)

const(
	flushTimeout = 5000
)
var errUnknownType = errors.New("unknown event type")

type Producer struct {
	producer *kafka.Producer
}

type TopUpEvent struct {
	WalletID  string    `json:"wallet_id"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// func NewProducer(address []string) (*Producer, error) {
func NewProducer(cnf config.Config) (*Producer, error) {
	kf_conf := &kafka.ConfigMap{
		"bootstrap:service": strings.Join(cnf.Kafka.Brokers, ","),
	}
	p, err := kafka.NewProducer(kf_conf)
	if err != nil {
		return nil, fmt.Errorf("error from producer", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: []byte(message),
		Key: nil,
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("not send message", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case *kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
