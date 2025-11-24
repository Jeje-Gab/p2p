package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

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
	Rarity     string `db:"rarity"`
	ImageURL   string `db:"image_url"`
}

// getRarityColor returns a gradient color based on rarity
func getRarityColor(rarity string) (string, string) {
	switch rarity {
	case "Consumer Grade":
		return "#B0C3D9", "#8FA9C4"
	case "Industrial Grade":
		return "#5E98D9", "#4B69FF"
	case "Mil-Spec Grade":
		return "#4B69FF", "#8847FF"
	case "Restricted":
		return "#8847FF", "#D32CE6"
	case "Classified":
		return "#D32CE6", "#EB4B4B"
	case "Covert":
		return "#EB4B4B", "#E4AE39"
	case "Exceedingly Rare":
		return "#E4AE39", "#FFD700"
	default:
		return "#CCCCCC", "#999999"
	}
}

// generateSVG creates a beautiful SVG placeholder for a skin
func generateSVG(skin Skin) string {
	color1, color2 := getRarityColor(skin.Rarity)

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="512" height="384" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="grad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" style="stop-color:%s;stop-opacity:1" />
      <stop offset="100%%" style="stop-color:%s;stop-opacity:1" />
    </linearGradient>
    <filter id="shadow">
      <feDropShadow dx="0" dy="4" stdDeviation="8" flood-opacity="0.3"/>
    </filter>
  </defs>

  <!-- Background gradient -->
  <rect width="512" height="384" fill="url(#grad)"/>

  <!-- Diagonal pattern -->
  <path d="M 0,384 L 512,0 L 512,100 L 0,484 Z" fill="white" opacity="0.05"/>
  <path d="M 0,284 L 512,-100 L 512,0 L 0,384 Z" fill="black" opacity="0.05"/>

  <!-- Center icon circle -->
  <circle cx="256" cy="142" r="60" fill="white" opacity="0.2"/>
  <circle cx="256" cy="142" r="50" fill="white" opacity="0.1"/>

  <!-- Weapon icon placeholder (crosshair) -->
  <circle cx="256" cy="142" r="40" fill="none" stroke="white" stroke-width="3" opacity="0.6"/>
  <line x1="256" y1="112" x2="256" y2="172" stroke="white" stroke-width="3" opacity="0.6"/>
  <line x1="226" y1="142" x2="286" y2="142" stroke="white" stroke-width="3" opacity="0.6"/>

  <!-- Weapon type badge -->
  <rect x="20" y="20" width="200" height="40" rx="20" fill="black" opacity="0.4"/>
  <text x="120" y="45" font-family="Arial, sans-serif" font-size="18" font-weight="bold" fill="white" text-anchor="middle" opacity="0.9">%s</text>

  <!-- Skin name -->
  <rect x="20" y="240" width="472" height="124" rx="12" fill="black" opacity="0.5" filter="url(#shadow)"/>
  <text x="256" y="285" font-family="Arial, sans-serif" font-size="32" font-weight="bold" fill="white" text-anchor="middle">%s</text>
  <text x="256" y="320" font-family="Arial, sans-serif" font-size="20" fill="white" text-anchor="middle" opacity="0.8">%s</text>

  <!-- Corner accent -->
  <path d="M 512,384 L 412,384 L 512,284 Z" fill="white" opacity="0.1"/>
</svg>`, color1, color2, skin.WeaponType, skin.Name, skin.Rarity)

	return svg
}

func main() {
	// Load .env file
	_ = godotenv.Load(".env")

	log.Println("[PLACEHOLDER] Starting placeholder image creation...")

	// Load configuration
	cfg := config.Load()

	// Initialize MinIO storage
	minioStorage, err := storage.NewMinIOStorage(cfg.MinIO)
	if err != nil {
		log.Fatalf("[PLACEHOLDER] Failed to initialize MinIO: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	sqlxdb, err := db.NewSqlx(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("[PLACEHOLDER] Failed to connect to database: %v", err)
	}
	defer sqlxdb.Close()

	// Get all skins
	skins, err := getAllSkins(sqlxdb)
	if err != nil {
		log.Fatalf("[PLACEHOLDER] Failed to get skins: %v", err)
	}

	log.Printf("[PLACEHOLDER] Found %d skins to process", len(skins))

	// Create placeholder for each skin
	successCount := 0
	errorCount := 0

	for i, skin := range skins {
		log.Printf("[PLACEHOLDER] [%d/%d] Processing: %s (%s)", i+1, len(skins), skin.Name, skin.WeaponType)

		// Generate object name
		objectName := generateObjectName(skin.Name, skin.WeaponType)

		// Check if already exists
		exists, err := minioStorage.FileExists(ctx, objectName)
		if err != nil {
			log.Printf("[PLACEHOLDER] Error checking file existence: %v", err)
			errorCount++
			continue
		}

		if exists {
			log.Printf("[PLACEHOLDER] File already exists, updating database: %s", objectName)
			newURL := minioStorage.GetPublicURL(objectName)
			if err := updateSkinImageURL(sqlxdb, skin.ID, newURL); err != nil {
				log.Printf("[PLACEHOLDER] Error updating database: %v", err)
				errorCount++
			} else {
				successCount++
			}
			continue
		}

		// Generate SVG
		svg := generateSVG(skin)

		// Upload to MinIO
		log.Printf("[PLACEHOLDER] Uploading to MinIO: %s", objectName)
		reader := bytes.NewReader([]byte(svg))
		newURL, err := minioStorage.UploadFile(ctx, objectName, reader, int64(len(svg)), "image/svg+xml")
		if err != nil {
			log.Printf("[PLACEHOLDER] Error uploading to MinIO: %v", err)
			errorCount++
			continue
		}

		// Update database
		log.Printf("[PLACEHOLDER] Updating database with new URL: %s", newURL)
		if err := updateSkinImageURL(sqlxdb, skin.ID, newURL); err != nil {
			log.Printf("[PLACEHOLDER] Error updating database: %v", err)
			errorCount++
			continue
		}

		successCount++
		log.Printf("[PLACEHOLDER] ✓ Success: %s -> %s", skin.Name, newURL)
	}

	log.Println("\n[PLACEHOLDER] ==================== SUMMARY ====================")
	log.Printf("[PLACEHOLDER] Total skins: %d", len(skins))
	log.Printf("[PLACEHOLDER] Successfully created: %d", successCount)
	log.Printf("[PLACEHOLDER] Errors: %d", errorCount)
	log.Println("[PLACEHOLDER] ==================================================")
}

func getAllSkins(db *sqlx.DB) ([]Skin, error) {
	var skins []Skin
	query := `SELECT id, name, weapon_type, rarity, image_url FROM p2p.skins WHERE deleted_at IS NULL`
	err := db.Select(&skins, query)
	return skins, err
}

func updateSkinImageURL(db *sqlx.DB, skinID int64, newURL string) error {
	query := `UPDATE p2p.skins SET image_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.Exec(query, newURL, skinID)
	return err
}

func generateObjectName(skinName, weaponType string) string {
	cleanName := ""
	for _, r := range skinName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cleanName += string(r)
		} else if r == ' ' || r == '-' {
			cleanName += "-"
		}
	}

	cleanWeapon := ""
	for _, r := range weaponType {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cleanWeapon += string(r)
		} else if r == ' ' || r == '-' {
			cleanWeapon += "-"
		}
	}

	return fmt.Sprintf("skins/%s/%s.svg", cleanWeapon, cleanName)
}
