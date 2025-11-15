-- Migration: Create offers and trades tables
-- For P2P skin trading functionality

-- Create offers table
CREATE TABLE IF NOT EXISTS p2p.offers (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES p2p.users(id) ON DELETE CASCADE,
    skin_offered_id BIGINT NOT NULL REFERENCES p2p.skins(id) ON DELETE CASCADE,
    skin_requested_id BIGINT NOT NULL REFERENCES p2p.skins(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'open', -- 'open', 'accepted', 'cancelled'
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT check_different_skins CHECK (skin_offered_id != skin_requested_id)
);

-- Index for listing open offers
CREATE INDEX idx_offers_status ON p2p.offers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_offers_owner ON p2p.offers(owner_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_offers_created ON p2p.offers(created_at DESC) WHERE deleted_at IS NULL;

-- Create trades table (completed trades)
CREATE TABLE IF NOT EXISTS p2p.trades (
    id BIGSERIAL PRIMARY KEY,
    offer_id BIGINT NOT NULL REFERENCES p2p.offers(id) ON DELETE CASCADE,
    from_user_id BIGINT NOT NULL REFERENCES p2p.users(id) ON DELETE CASCADE,
    to_user_id BIGINT NOT NULL REFERENCES p2p.users(id) ON DELETE CASCADE,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_different_users CHECK (from_user_id != to_user_id)
);

-- Index for trade history queries
CREATE INDEX idx_trades_from_user ON p2p.trades(from_user_id);
CREATE INDEX idx_trades_to_user ON p2p.trades(to_user_id);
CREATE INDEX idx_trades_executed ON p2p.trades(executed_at DESC);
CREATE INDEX idx_trades_offer_id ON p2p.trades(offer_id);

COMMENT ON TABLE p2p.offers IS 'P2P trade offers created by users';
COMMENT ON COLUMN p2p.offers.status IS 'Offer lifecycle: open (available) -> accepted (trade executed) or cancelled (user cancelled)';
COMMENT ON TABLE p2p.trades IS 'Historical record of completed trades for audit trail';
