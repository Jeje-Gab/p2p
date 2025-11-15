package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/cs2-p2p-skins/backend/pkg/auth"
	"github.com/cs2-p2p-skins/backend/pkg/config"
	"github.com/cs2-p2p-skins/backend/pkg/db"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"

	authHttp "github.com/cs2-p2p-skins/backend/internal/auth/delivery/http"
	authRepo "github.com/cs2-p2p-skins/backend/internal/auth/repository"
	authUC "github.com/cs2-p2p-skins/backend/internal/auth/usecase"

	skinsHttp "github.com/cs2-p2p-skins/backend/internal/skins/delivery/http"
	skinsRepo "github.com/cs2-p2p-skins/backend/internal/skins/repository"
	skinsUC "github.com/cs2-p2p-skins/backend/internal/skins/usecase"

	offersHttp "github.com/cs2-p2p-skins/backend/internal/offers/delivery/http"
	offersRepo "github.com/cs2-p2p-skins/backend/internal/offers/repository"
	offersUC "github.com/cs2-p2p-skins/backend/internal/offers/usecase"

	tradesHttp "github.com/cs2-p2p-skins/backend/internal/trades/delivery/http"
	tradesRepo "github.com/cs2-p2p-skins/backend/internal/trades/repository"
	tradesUC "github.com/cs2-p2p-skins/backend/internal/trades/usecase"
)

func main() {
	// Load .env file
	_ = godotenv.Load(".env")

	// Setup graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	sqlxdb, err := db.NewSqlx(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer sqlxdb.Close()

	log.Printf("[BOOT] Environment: %s | Port: %s", cfg.Env, cfg.HTTP.Port)

	// Initialize auth managers
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
	totpManager := auth.NewTOTPManager("CS2 P2P Skins")
	steamAuth := auth.NewSteamOAuthManager(cfg.Steam.APIKey, cfg.Steam.CallbackURL)

	// Initialize repositories
	authRepository := authRepo.New(sqlxdb)
	skinsRepository := skinsRepo.New(sqlxdb)
	offersRepository := offersRepo.New(sqlxdb)
	tradesRepository := tradesRepo.New(sqlxdb)

	// Initialize use cases
	authUseCase := authUC.New(authRepository, jwtManager, totpManager, steamAuth)
	skinsUseCase := skinsUC.New(skinsRepository)
	offersUseCase := offersUC.New(offersRepository, skinsRepository, tradesRepository)
	tradesUseCase := tradesUC.New(tradesRepository)

	// Initialize rate limiters
	strictRL := middleware.StrictRateLimiter()  // For auth endpoints
	moderateRL := middleware.ModerateRateLimiter() // For API endpoints

	// Initialize Echo server
	e := echo.New()
	e.HideBanner = true

	// Global middlewares
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.SecurityHeaders())
	e.Use(echoMiddleware.CORSWithConfig(middleware.CORSConfig(
		strings.Split(cfg.AllowOrigins, ","),
	)))

	// Health check endpoint
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := e.Group("/api")

	// Auth routes (public + protected)
	authGroup := api.Group("/auth")
	authHttp.RegisterRoutes(authGroup, authUseCase, jwtManager, strictRL)

	// Protected routes - require JWT authentication
	protected := api.Group("", middleware.JWTAuth(jwtManager), moderateRL.Middleware())

	// Skins routes
	skinsGroup := protected.Group("/skins")
	skinsHttp.RegisterRoutes(skinsGroup, skinsUseCase)

	// Offers routes
	offersGroup := protected.Group("/offers")
	offersHttp.RegisterRoutes(offersGroup, offersUseCase)

	// Trades routes
	tradesGroup := protected.Group("/trades")
	tradesHttp.RegisterRoutes(tradesGroup, tradesUseCase)

	// Start server in goroutine
	go func() {
		log.Printf("[HTTP] Server starting on port %s", cfg.HTTP.Port)
		if err := e.Start(":" + cfg.HTTP.Port); err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Graceful shutdown
	log.Println("[SHUTDOWN] Shutting down server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown error:", err)
	}

	log.Println("[SHUTDOWN] Server stopped")
}
