package user

import (
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/repositories/user"
	"errors"
	"log"
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
		return nil, errors.New("username or Password or Email or Phone is empty")
	}

	if len(user.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	return u.r.Add(user)
}

func (u *UserUseCases) FindById(userid int64) (*models.User, error) {
	return u.r.FindById(userid)
}

func (u *UserUseCases) Delete(userId int64) error {
	if res, err := u.FindById(userId); err != nil && res != nil {
		err = u.r.Delete(userId)

		if err != nil {
			log.Printf("Error deleting user: %v", err)
		}
	}

	return errors.New("user not found")
}
