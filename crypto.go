package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	// this is a good decent default, if we need to go higher,
	// our clients are into some shady shit and they deserve what they get.

	// BCryptRounds the default number of rounds passed to bcrypt.
	BCryptRounds = 11
)

// HashPassword hashes a password using bcrypt and returns the string representation of it.
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", ErrNoPassword
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), BCryptRounds)
	return string(h), err
}

// CheckPassword checks a hashed password against a plain-text password.
func CheckPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// IsHashedPass checks if a password hash is a valid bcrypt hash or not.
func IsHashedPass(hash string) bool {
	if hash == "" {
		return false
	}
	cost, err := bcrypt.Cost([]byte(hash))
	return err == nil && cost >= BCryptRounds
}

// RandomToken returns a random `string` crypto/rand generated token with the given length.
// If b64 is true, it will encode it with base64.RawURLEncoding otherwise uses hex.
func RandomToken(ln int, b64 bool) string {
	tok := make([]byte, ln)
	if n, _ := rand.Read(tok); n != len(tok) {
		log.Panicf("expected %d rand bytes, got %d, something is wrong", len(tok), n)
	}
	if b64 {
		return base64.RawURLEncoding.EncodeToString(tok)
	}
	return hex.EncodeToString(tok)
}
