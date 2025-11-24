package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cs2-p2p-skins/backend/pkg/config"
	"github.com/cs2-p2p-skins/backend/pkg/db"
	"github.com/cs2-p2p-skins/backend/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Skin struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	WeaponType string `db:"weapon_type"`
	ImageURL   string `db:"image_url"`
}

func main() {
	// Load .env file
	_ = godotenv.Load(".env")

	log.Println("[MIGRATE] Starting image migration to MinIO...")

	// Load configuration
	cfg := config.Load()

	// Initialize MinIO storage
	minioStorage, err := storage.NewMinIOStorage(cfg.MinIO)
	if err != nil {
		log.Fatalf("[MIGRATE] Failed to initialize MinIO: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	sqlxdb, err := db.NewSqlx(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("[MIGRATE] Failed to connect to database: %v", err)
	}
	defer sqlxdb.Close()

	// Get all skins
	skins, err := getAllSkins(sqlxdb)
	if err != nil {
		log.Fatalf("[MIGRATE] Failed to get skins: %v", err)
	}

	log.Printf("[MIGRATE] Found %d skins to migrate", len(skins))

	// Migrate each skin
	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, skin := range skins {
		log.Printf("[MIGRATE] [%d/%d] Processing: %s (%s)", i+1, len(skins), skin.Name, skin.WeaponType)

		// Skip if not a Steam URL
		if !strings.Contains(skin.ImageURL, "steamstatic.com") {
			log.Printf("[MIGRATE] Skipping (not a Steam URL): %s", skin.Name)
			skipCount++
			continue
		}

		// Generate object name
		objectName := generateObjectName(skin.Name, skin.WeaponType)

		// Check if already exists
		exists, err := minioStorage.FileExists(ctx, objectName)
		if err != nil {
			log.Printf("[MIGRATE] Error checking file existence: %v", err)
			errorCount++
			continue
		}

		if exists {
			log.Printf("[MIGRATE] File already exists, updating database: %s", objectName)
			newURL := minioStorage.GetPublicURL(objectName)
			if err := updateSkinImageURL(sqlxdb, skin.ID, newURL); err != nil {
				log.Printf("[MIGRATE] Error updating database: %v", err)
				errorCount++
			} else {
				successCount++
			}
			continue
		}

		// Download image from Steam
		log.Printf("[MIGRATE] Downloading from Steam: %s", skin.ImageURL)
		resp, err := http.Get(skin.ImageURL)
		if err != nil {
			log.Printf("[MIGRATE] Error downloading image: %v", err)
			errorCount++
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Printf("[MIGRATE] HTTP error %d downloading image", resp.StatusCode)
			errorCount++
			continue
		}

		// Upload to MinIO
		log.Printf("[MIGRATE] Uploading to MinIO: %s", objectName)
		newURL, err := minioStorage.UploadFile(ctx, objectName, resp.Body, resp.ContentLength, "image/png")
		resp.Body.Close()

		if err != nil {
			log.Printf("[MIGRATE] Error uploading to MinIO: %v", err)
			errorCount++
			continue
		}

		// Update database
		log.Printf("[MIGRATE] Updating database with new URL: %s", newURL)
		if err := updateSkinImageURL(sqlxdb, skin.ID, newURL); err != nil {
			log.Printf("[MIGRATE] Error updating database: %v", err)
			errorCount++
			continue
		}

		successCount++
		log.Printf("[MIGRATE] ✓ Success: %s -> %s", skin.Name, newURL)
	}

	log.Println("\n[MIGRATE] ==================== SUMMARY ====================")
	log.Printf("[MIGRATE] Total skins: %d", len(skins))
	log.Printf("[MIGRATE] Successfully migrated: %d", successCount)
	log.Printf("[MIGRATE] Skipped: %d", skipCount)
	log.Printf("[MIGRATE] Errors: %d", errorCount)
	log.Println("[MIGRATE] ===================================================")
}

func getAllSkins(db *sqlx.DB) ([]Skin, error) {
	var skins []Skin
	query := `SELECT id, name, weapon_type, image_url FROM p2p.skins WHERE deleted_at IS NULL`
	err := db.Select(&skins, query)
	return skins, err
}

func updateSkinImageURL(db *sqlx.DB, skinID int64, newURL string) error {
	query := `UPDATE p2p.skins SET image_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.Exec(query, newURL, skinID)
	return err
}

func generateObjectName(skinName, weaponType string) string {
	// Clean the names
	cleanName := strings.ToLower(strings.ReplaceAll(skinName, " ", "-"))
	cleanWeapon := strings.ToLower(strings.ReplaceAll(weaponType, " ", "-"))

	// Remove special characters
	cleanName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, cleanName)

	cleanWeapon = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, cleanWeapon)

	return fmt.Sprintf("skins/%s/%s.png", cleanWeapon, cleanName)
}
