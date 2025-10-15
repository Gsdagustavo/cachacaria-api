package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var ErrExpiredToken = errors.New("TOKEN_HAS_EXPIRED")

// Remade Maker interface
type Maker interface {
	// UPDATED: Now requires the userID (int) when creating a token
	CreateToken(email string, userID int, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

// Remade Payload struct
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserID    int       `json:"user_id"` // NEW FIELD for context storage
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

// Remade NewPayload function
func NewPayload(email string, userID int, duration time.Duration) (*Payload, error) { // UPDATED signature
	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token uuid: %s", err)
	}

	payload := &Payload{
		ID:        tokenUUID,
		Email:     email,
		UserID:    userID, // SET the new field
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Concrete struct definition for PasetoMaker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

func NewPasetoMaker(symmetricKey string) Maker {
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}
}

// Remade PasetoMaker.CreateToken method
func (m *PasetoMaker) CreateToken(email string, userID int, duration time.Duration) (string, error) { // UPDATED signature
	payload, err := NewPayload(email, userID, duration) // UPDATED call
	if err != nil {
		return "", err
	}

	encrypted, err := m.paseto.Encrypt([]byte(m.symmetricKey), payload, nil)
	if err != nil {
		return "", fmt.Errorf("error encrypting token: %s", err)
	}

	return encrypted, nil
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := m.paseto.Decrypt(token, []byte(m.symmetricKey), &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting token: %s", err)
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
