package user

import (
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/repositories/user"
	"errors"
)

type UserUseCases struct {
	r *user.UserRepository
}

func NewUserUseCases(r *user.UserRepository) *UserUseCases {
	return &UserUseCases{r}
}

func (u *UserUseCases) GetAll() []models.User {
	return u.r.GetAll()
}

func (u *UserUseCases) Add(user models.AddUserRequest) (*models.AddUserResponse, error) {
	if user.Name == "" || user.Password == "" || user.Email == "" || user.Phone == "" {
		return nil, errors.New("Username or Password or Email or Phone is empty")
	}

	if len(user.Password) < 8 {
		return nil, errors.New("Password must be at least 8 characters")
	}

	return u.r.Add(user)
}
