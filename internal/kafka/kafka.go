package kafka

import (
	"encoding/json"
	"time"
	"wallet-api-service/internal/config"
	// "wallet-api-service/internal/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Client struct {
	cfg *config.Config
}

type TopUpEvent struct {
	WalletID  string    `json:"wallet_id"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func New(cfg config.Config) *Client {
	return &Client{
		cfg: &cfg,
	}
}

func (c *Client) PublishTopUp(walletID uuid.UUID, amount int) error {
	event := TopUpEvent{
		WalletID:  walletID.String(),
		Amount:    amount,
		Timestamp: time.Now(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// В реальной имплементации здесь был бы код для отправки сообщения в Kafka
	// producer.SendMessage(...)
	
	log.Info().
		Str("wallet_id", walletID.String()).
		Int("amount", amount).
		Msg("Published top-up event to Kafka")
	
	log.Debug().RawJSON("event", eventJSON).Msg("Raw Kafka event")
	
	return nil
}

func (c *Client) Close() error {
	return nil
}
