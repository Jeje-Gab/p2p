package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the bcrypt cost for password hashing
	// Higher cost = more secure but slower
	// 12 is a good balance between security and performance
	DefaultCost = 12
)

// HashPassword generates a bcrypt hash from a plain password
// Security: NEVER store passwords in plain text!
// Bcrypt automatically handles salt generation
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if a plain password matches the stored hash
// Returns true if password is correct, false otherwise
// Security: Constant-time comparison prevents timing attacks
func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
