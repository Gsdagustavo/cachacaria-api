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
	_, err := u.FindById(userId)
	if err != nil {
		return errors.New("user not found. error: " + err.Error())
	}

	err = u.r.Delete(userId)
	if err != nil {
		return errors.New("user could not be deleted. error: " + err.Error())
	}

	return nil
}
