package jwt

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher struct for hashing and verifying passwords
type PasswordHasher struct {
	cost int // Cost parameter for bcrypt
}

// NewPasswordHasher creates a new PasswordHasher with the given cost
func NewPasswordHasher(cost int) PasswordHasher {
	return PasswordHasher{cost: cost}
}

// HashPassword hashes the given plain-text password
func (ph PasswordHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a plain-text password with a hashed password
func (ph PasswordHasher) VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
