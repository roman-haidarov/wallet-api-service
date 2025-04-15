package kafka

import (
	// "encoding/json"
	"fmt"
	"strings"
	"time"
	// "wallet-api-service/internal/config"

	// "wallet-api-service/internal/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/google/uuid"
	// "github.com/rs/zerolog/log"
)

type Producer struct {
	producer *kafka.Producer
}

type TopUpEvent struct {
	WalletID  string    `json:"wallet_id"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap:service": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error from producer", err)
	}

	return &Producer{producer: p}, nil
}

// type Client struct {
// 	cfg *config.Config
// }

// func New(cfg config.Config) *Client {
// 	return &Client{
// 		cfg: &cfg,
// 	}
// }

// func (c *Client) PublishTopUp(walletID uuid.UUID, amount int) error {
// 	event := TopUpEvent{
// 		WalletID:  walletID.String(),
// 		Amount:    amount,
// 		Timestamp: time.Now(),
// 	}

// 	eventJSON, err := json.Marshal(event)
// 	if err != nil {
// 		return err
// 	}

// 	// В реальной имплементации здесь был бы код для отправки сообщения в Kafka
// 	// producer.SendMessage(...)
	
// 	log.Info().
// 		Str("wallet_id", walletID.String()).
// 		Int("amount", amount).
// 		Msg("Published top-up event to Kafka")
	
// 	log.Debug().RawJSON("event", eventJSON).Msg("Raw Kafka event")
	
// 	return nil
// }

// func (c *Client) Close() error {
// 	return nil
// }
