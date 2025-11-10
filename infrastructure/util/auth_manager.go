package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrExpiredToken = errors.New("TOKEN_HAS_EXPIRED")
	ErrInvalidToken = errors.New("INVALID_TOKEN")
)

const (
	DefaultCost          = bcrypt.DefaultCost
	DefaultTokenDuration = 12 * time.Hour
)

type AuthManager struct {
	paseto       *paseto.V2
	symmetricKey string
}

// NewAuthManager creates a new manager instance with a given symmetric key.
func NewAuthManager(symmetricKey string) AuthManager {
	return AuthManager{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserID    int       `json:"user_id"`
	IsAdmin   bool      `json:"is_admin"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(email string, userID int, isAdmin bool, duration time.Duration) (*Payload, error) {
	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token uuid: %w", err)
	}

	return &Payload{
		ID:        tokenUUID,
		Email:     email,
		UserID:    userID,
		IsAdmin:   isAdmin,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

func (a *AuthManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(bytes), nil
}

func (a *AuthManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthManager) CreateToken(email string, userID int, isAdmin bool) (string, error) {
	payload, err := NewPayload(email, userID, isAdmin, DefaultTokenDuration)
	if err != nil {
		return "", err
	}

	encrypted, err := a.paseto.Encrypt([]byte(a.symmetricKey), payload, nil)
	if err != nil {
		return "", fmt.Errorf("error encrypting token: %w", err)
	}
	return encrypted, nil
}

func (a *AuthManager) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := a.paseto.Decrypt(token, []byte(a.symmetricKey), &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}

	return &payload, nil
}
