package types

import (
	"time"
	"github.com/google/uuid"
)

type Wallet struct {
	WalletID  uuid.UUID `json:"wallet_id" db:"wallet_id"`
	Amount    int       `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
