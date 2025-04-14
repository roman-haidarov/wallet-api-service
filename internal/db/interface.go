package db

import (
	"context"
	"wallet-api-service/internal/types"
	"github.com/google/uuid"
)

type DB interface {
	GetWallet(ctx context.Context, walletID uuid.UUID) (*types.Wallet, error)
	TopUpWallet(ctx context.Context, walletID uuid.UUID, amount int) error
	CreateWallet(ctx context.Context, wallet types.Wallet) error
}
