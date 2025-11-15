-- Migration: Create users and login_attempts tables
-- Security features: password hashing, 2FA, role-based access control, brute force protection

-- Create users table with security features
CREATE TABLE IF NOT EXISTS p2p.users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL, -- bcrypt hash, never store plain passwords!
    steam_id VARCHAR(100) UNIQUE, -- For Steam OAuth integration
    twofa_secret VARCHAR(255), -- TOTP secret for 2FA (Google Authenticator)
    twofa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- 'user' or 'admin' for authorization
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP -- Soft delete for data retention
);

-- Index for faster lookups
CREATE INDEX idx_users_email ON p2p.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_steam_id ON p2p.users(steam_id) WHERE steam_id IS NOT NULL AND deleted_at IS NULL;

-- Create login_attempts table for brute force protection
CREATE TABLE IF NOT EXISTS p2p.login_attempts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES p2p.users(id) ON DELETE SET NULL, -- Nullable if user doesn't exist
    ip_address VARCHAR(45) NOT NULL, -- IPv4 or IPv6
    success BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for brute force detection queries
CREATE INDEX idx_login_attempts_user_created ON p2p.login_attempts(user_id, created_at DESC);
CREATE INDEX idx_login_attempts_ip_created ON p2p.login_attempts(ip_address, created_at DESC);

-- Comments explaining security decisions
COMMENT ON TABLE p2p.users IS 'User accounts with password hashing (bcrypt/argon2), 2FA support, and role-based access';
COMMENT ON COLUMN p2p.users.password_hash IS 'Bcrypt or Argon2id hash - NEVER store passwords in plain text!';
COMMENT ON COLUMN p2p.users.twofa_secret IS 'TOTP secret for Google Authenticator - provides second factor authentication';
COMMENT ON TABLE p2p.login_attempts IS 'Tracks login attempts for brute force detection and mitigation';
