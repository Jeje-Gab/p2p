-- Migration: Create backup_codes table for 2FA recovery
-- Security: Provides one-time use recovery codes if user loses access to TOTP device

-- Create backup_codes table
CREATE TABLE IF NOT EXISTS p2p.backup_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES p2p.users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL, -- Bcrypt hash of backup code, never store plain!
    used BOOLEAN NOT NULL DEFAULT FALSE, -- Track if code has been used
    used_at TIMESTAMP, -- When the code was used
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for faster lookups
CREATE INDEX idx_backup_codes_user_id ON p2p.backup_codes(user_id) WHERE used = FALSE;
CREATE INDEX idx_backup_codes_user_unused ON p2p.backup_codes(user_id, used);

-- Comments explaining security decisions
COMMENT ON TABLE p2p.backup_codes IS 'One-time use backup codes for 2FA recovery when user loses access to authenticator app';
COMMENT ON COLUMN p2p.backup_codes.code_hash IS 'Bcrypt hash of 8-character backup code - never store plain text!';
COMMENT ON COLUMN p2p.backup_codes.used IS 'Once used, backup code cannot be reused (prevents replay attacks)';
