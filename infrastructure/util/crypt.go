package util

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = bcrypt.DefaultCost
const DefaultTokenDuration = 12 * time.Hour // hours

// Remade Crypt interface
type Crypt interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool

	// UPDATED: Now requires the userID (int)
	GenerateAuthToken(email string, userID int) (string, error)

	// NEW: Exposes token verification to use cases/middleware
	VerifyAuthToken(token string) (*Payload, error)
}

type crypt struct {
	maker Maker
}

func NewCrypt(maker Maker) Crypt {
	return &crypt{
		maker: maker,
	}
}

func (c crypt) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c crypt) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Remade crypt.GenerateAuthToken method
func (c crypt) GenerateAuthToken(email string, userID int) (string, error) { // Corrected receiver to concrete type 'crypt'
	// Delegation to the Maker
	return c.maker.CreateToken(email, userID, DefaultTokenDuration)
}

// NEW: crypt.VerifyAuthToken method
func (c crypt) VerifyAuthToken(token string) (*Payload, error) { // Corrected receiver to concrete type 'crypt'
	// Delegation to the Maker
	return c.maker.VerifyToken(token)
}
