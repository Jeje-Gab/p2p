-- Migration: Create skins and user_skins tables
-- For CS2 skin inventory management

-- Create skins table (master list of all available skins)
CREATE TABLE IF NOT EXISTS skins (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    weapon_type VARCHAR(100) NOT NULL, -- e.g., 'AK-47', 'AWP', 'M4A4'
    rarity VARCHAR(50) NOT NULL, -- Consumer Grade, Industrial Grade, Mil-Spec, Restricted, Classified, Covert
    float_value DECIMAL(10, 8), -- Wear value from 0.0 (Factory New) to 1.0 (Battle-Scarred)
    image_url TEXT NOT NULL, -- URL to skin image
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Index for searching skins
CREATE INDEX idx_skins_weapon_type ON skins(weapon_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_skins_rarity ON skins(rarity) WHERE deleted_at IS NULL;

-- Create user_skins table (junction table for user inventory)
CREATE TABLE IF NOT EXISTS user_skins (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skin_id BIGINT NOT NULL REFERENCES skins(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, skin_id) -- A user can't have duplicate entries for the same skin
);

-- Index for faster inventory queries
CREATE INDEX idx_user_skins_user_id ON user_skins(user_id);
CREATE INDEX idx_user_skins_skin_id ON user_skins(skin_id);

COMMENT ON TABLE skins IS 'Master list of CS2 weapon skins available in the system';
COMMENT ON TABLE user_skins IS 'User inventory - tracks which skins each user owns and quantity';
COMMENT ON COLUMN skins.float_value IS 'Skin wear level: 0.00-0.07 Factory New, 0.07-0.15 Minimal Wear, 0.15-0.38 Field-Tested, 0.38-0.45 Well-Worn, 0.45-1.00 Battle-Scarred';
