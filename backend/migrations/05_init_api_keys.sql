-- Migration 05: API Keys for External Access
-- Creates table for managing external API keys

-- API Keys table
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    key_prefix VARCHAR(30) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);

-- Comments
COMMENT ON TABLE api_keys IS 'External API keys for third-party integrations';
COMMENT ON COLUMN api_keys.key_hash IS 'SHA256 hash of the API key (for security)';
COMMENT ON COLUMN api_keys.key_prefix IS 'First 20 characters of key for display (e.g., "sk_live_abc123456789...")';
COMMENT ON COLUMN api_keys.name IS 'Friendly name for the API key (e.g., "Mobile App", "Analytics Dashboard")';
COMMENT ON COLUMN api_keys.description IS 'Optional description of key usage';
COMMENT ON COLUMN api_keys.is_active IS 'Whether the key is active (can be revoked)';
COMMENT ON COLUMN api_keys.last_used_at IS 'Last time this key was used to authenticate';
COMMENT ON COLUMN api_keys.expires_at IS 'Optional expiration date for the key';
