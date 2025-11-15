package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env          string
	HTTP         HTTPConfig
	DB           DBConfig
	JWT          JWTConfig
	Steam        SteamConfig
	AllowOrigins string
}

type HTTPConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type SteamConfig struct {
	APIKey      string
	CallbackURL string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	// JWT expiration in hours (default 24h)
	jwtExpHours, _ := strconv.Atoi(getenv("JWT_EXPIRATION_HOURS", "24"))

	return &Config{
		Env: getenv("APP_ENV", "dev"),
		HTTP: HTTPConfig{
			Port: getenv("PORT", "8080"),
		},
		DB: DBConfig{
			DSN: getenv("DB_DSN", "postgres://postgres:1212@localhost:5432/postgres?sslmode=disable"),
		},
		JWT: JWTConfig{
			Secret:     getenv("JWT_SECRET", "change-me-in-production"),
			Expiration: time.Duration(jwtExpHours) * time.Hour,
		},
		Steam: SteamConfig{
			APIKey:      getenv("STEAM_API_KEY", ""),
			CallbackURL: getenv("STEAM_CALLBACK_URL", "http://localhost:8080/api/auth/steam/callback"),
		},
		AllowOrigins: getenv("ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}
