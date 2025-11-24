# ==========================================
# Multi-stage Dockerfile for Go Backend
# ==========================================

# Stage 1: Build
FROM golang:1.24 AS build

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for sqlite/postgres drivers
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /app/server ./cmd/server

# Install golang-migrate for database migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ==========================================
# Stage 2: Runtime
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    postgresql-client \
    curl \
    bash

WORKDIR /app

# Copy binary from builder
COPY --from=build /app/server /app/server
COPY --from=build /go/bin/migrate /usr/local/bin/migrate

# Copy migrations and scripts
COPY migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

# Make entrypoint executable
RUN chmod +x /app/docker-entrypoint.sh

# Environment variables
ENV APP_ENV=production
ENV PORT=8080
ENV TZ=America/Sao_Paulo

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

# Run entrypoint script
ENTRYPOINT ["/app/docker-entrypoint.sh"]
