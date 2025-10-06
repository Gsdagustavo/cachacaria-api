package entities

import "github.com/google/uuid"

type User struct {
	ID       int64     `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Phone    string    `json:"phone"`
	IsAdm    bool      `json:"is_adm"`
}

type UserCredentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	IsAdm    bool   `json:"is_adm"`
}
