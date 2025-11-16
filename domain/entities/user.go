package entities

import "github.com/google/uuid"

type User struct {
	ID       int       `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Phone    string    `json:"phone"`
	IsAdm    bool      `json:"is_adm"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	IsAdm    bool   `json:"is_adm,omitempty"`
}

type ChangePasswordRequest struct {
	UserID                  int64  `json:"user_id"`
	CurrentPassword         string `json:"password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}
