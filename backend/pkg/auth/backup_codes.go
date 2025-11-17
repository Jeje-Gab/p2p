package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

const (
	// Number of backup codes to generate
	BackupCodesCount = 10
	// Length of each backup code (8 characters)
	BackupCodeLength = 8
)

// BackupCode represents a single backup code
type BackupCode struct {
	PlainCode string // Plain text code (show to user once)
	Hash      string // Bcrypt hash (store in database)
}

// GenerateBackupCodes creates N random backup codes
// Returns plain codes (to show user) and hashed codes (to store in DB)
// Security: Codes are cryptographically random and stored hashed
func GenerateBackupCodes(count int) ([]BackupCode, error) {
	codes := make([]BackupCode, count)

	for i := 0; i < count; i++ {
		// Generate random bytes
		randomBytes := make([]byte, 10) // 10 bytes = 16 base32 chars
		if _, err := rand.Read(randomBytes); err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}

		// Encode to base32 (uppercase, no padding)
		plainCode := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
		// Take first 8 characters and format as XXXX-XXXX
		plainCode = strings.ToUpper(plainCode[:BackupCodeLength])
		formattedCode := fmt.Sprintf("%s-%s", plainCode[:4], plainCode[4:])

		// Hash the code for storage
		hash, err := HashPassword(formattedCode)
		if err != nil {
			return nil, fmt.Errorf("failed to hash backup code: %w", err)
		}

		codes[i] = BackupCode{
			PlainCode: formattedCode,
			Hash:      hash,
		}
	}

	return codes, nil
}

// ValidateBackupCode verifies if a backup code matches the stored hash
// Security: Uses bcrypt's constant-time comparison
func ValidateBackupCode(hash, code string) bool {
	// Normalize code format (remove spaces, uppercase)
	code = strings.ToUpper(strings.ReplaceAll(code, " ", ""))
	return VerifyPassword(hash, code)
}
