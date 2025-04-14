package memdb

import (
	"context"
	"sync"
	"time"
	"errors"
	"wallet-api-service/internal/types"
	"wallet-api-service/internal/db"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type MemDB struct {
	mu       sync.Mutex
	wallets  map[uuid.UUID]types.Wallet
}

func New() db.DB {
	return &MemDB{
		wallets: make(map[uuid.UUID]types.Wallet),
	}
}

func (db *MemDB) GetWallet(ctx context.Context, walletID uuid.UUID) (*types.Wallet, error) {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "DB:GetWallet")
	defer span.End()

	db.mu.Lock()
	defer db.mu.Unlock()

	wallet, exists := db.wallets[walletID]
	if !exists {
		return nil, errors.New("wallet not found")
	}

	return &wallet, nil
}

func (db *MemDB) TopUpWallet(ctx context.Context, walletID uuid.UUID, amount int) error {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "DB:TopUpWallet")
	defer span.End()

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	wallet, exists := db.wallets[walletID]
	if !exists {
		return errors.New("wallet not found")
	}

	wallet.Amount += amount
	wallet.UpdatedAt = time.Now()
	db.wallets[walletID] = wallet

	return nil
}

func (db *MemDB) CreateWallet(ctx context.Context, wallet types.Wallet) error {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "DB:CreateWallet")
	defer span.End()

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.wallets[wallet.WalletID]; exists {
		return errors.New("wallet already exists")
	}

	db.wallets[wallet.WalletID] = wallet
	return nil
}
