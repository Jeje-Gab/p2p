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
	TLS          TLSConfig
	MinIO        MinIOConfig
	APIKeys      string
	AllowOrigins string
}

type HTTPConfig struct {
	Port string
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
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

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	PublicURL string
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
		TLS: TLSConfig{
			Enabled:  getenv("TLS_ENABLED", "false") == "true",
			CertFile: getenv("TLS_CERT_FILE", "certs/server.crt"),
			KeyFile:  getenv("TLS_KEY_FILE", "certs/server.key"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getenv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getenv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getenv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getenv("MINIO_BUCKET", "net.public.p2p"),
			UseSSL:    getenv("MINIO_USE_SSL", "false") == "true",
			PublicURL: getenv("MINIO_PUBLIC_URL", "http://localhost:9000"),
		},
		APIKeys:      getenv("EXTERNAL_API_KEYS", ""),
		AllowOrigins: getenv("ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}
