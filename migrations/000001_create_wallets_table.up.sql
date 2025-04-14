CREATE TABLE IF NOT EXISTS wallets (
    wallet_id UUID PRIMARY KEY,
    amount INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_wallets_created_at ON wallets(created_at);
