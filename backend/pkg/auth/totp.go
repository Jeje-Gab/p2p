package auth

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

// TOTPManager handles Time-based One-Time Password (2FA)
type TOTPManager struct {
	issuer string
}

// NewTOTPManager creates a new TOTP manager
func NewTOTPManager(issuer string) *TOTPManager {
	return &TOTPManager{
		issuer: issuer,
	}
}

// GenerateSecret creates a new TOTP secret for a user
// Returns the secret string and a QR code URL for Google Authenticator
// Security: The secret should be stored encrypted and never exposed in logs
func (m *TOTPManager) GenerateSecret(email string) (secret string, qrURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      m.issuer,
		AccountName: email,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	return key.Secret(), key.URL(), nil
}

// ValidateCode verifies if a TOTP code is valid for the given secret
// Security: Codes are time-based (30-second window) and can only be used once
func (m *TOTPManager) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
