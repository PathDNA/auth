package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"

	"github.com/missionMeteora/toolkit/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrInvalidToken = errors.Error("invalid token")
	ErrMissingId    = errors.Error("missing id")
	ErrInvalidLogin = errors.Error("invalid login")

	// this is a good decent default, if we need to go higher,
	// our clients are into some shady shit and they deserve what they get.
	bcryptRounds = 11
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", nil
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptRounds)
	return string(h), err
}

func CheckPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func IsHashedPass(hash string) bool {
	if hash == "" {
		return false
	}
	cost, err := bcrypt.Cost([]byte(hash))
	return err == nil && cost >= bcryptRounds
}

func CreateMAC(password, token, salt string) string {
	h := hmac.New(sha256.New, []byte(token+salt))
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyMac(mac1, password, token, salt string) bool {
	mac2 := decodeHex(CreateMAC(password, token, salt))
	return hmac.Equal(decodeHex(mac1), mac2)
}

// RandomToken returns a `string` crypto/rand generated token with the given length.
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

func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}
