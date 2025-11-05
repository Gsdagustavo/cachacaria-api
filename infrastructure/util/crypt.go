package util

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = bcrypt.DefaultCost
const DefaultTokenDuration = 12 * time.Hour // hours

type Crypt interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateAuthToken(email string, userID int, isAdmin bool) (string, error)
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

func (c crypt) GenerateAuthToken(email string, userID int, isAdmin bool) (string, error) { // Corrected receiver to concrete type 'crypt'
	return c.maker.CreateToken(email, userID, isAdmin, DefaultTokenDuration)
}

func (c crypt) VerifyAuthToken(token string) (*Payload, error) { // Corrected receiver to concrete type 'crypt'
	return c.maker.VerifyToken(token)
}
