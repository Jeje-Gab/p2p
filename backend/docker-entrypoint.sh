#!/bin/bash
set -e

echo "=========================================="
echo "CS2 P2P Skins - Backend Startup"
echo "=========================================="

# Wait for PostgreSQL to be ready
echo "[INIT] Waiting for PostgreSQL..."
until pg_isready -h "${DB_HOST:-postgres}" -p "${DB_PORT:-5432}" -U "${DB_USER:-postgres}"; do
  echo "[INIT] PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "[INIT] PostgreSQL is ready!"

# Run migrations
if [ -d "/app/migrations" ]; then
  echo "[MIGRATIONS] Running database migrations..."

  # Build the migration database URL
  DB_URL="postgres://${DB_USER:-postgres}:${DB_PASSWORD:-postgres}@${DB_HOST:-postgres}:${DB_PORT:-5432}/${DB_NAME:-postgres}?sslmode=disable"

  # Run migrations
  migrate -path=/app/migrations -database="${DB_URL}" up

  if [ $? -eq 0 ]; then
    echo "[MIGRATIONS] Migrations completed successfully!"
  else
    echo "[MIGRATIONS] Warning: Some migrations may have failed or were already applied"
  fi
else
  echo "[MIGRATIONS] No migrations directory found, skipping..."
fi

echo "[INIT] Starting application..."
echo "=========================================="

# Execute the main application
exec /app/server
