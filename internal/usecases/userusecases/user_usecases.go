package userusecases

import (
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/repositories/userrepository"
	"errors"
)

type UserUseCases struct {
	r *userrepository.UserRepository
}

func NewUserUseCases(r *userrepository.UserRepository) *UserUseCases {
	return &UserUseCases{r}
}

func (u *UserUseCases) GetAll() ([]models.User, error) {
	return u.r.GetAll()
}

func (u *UserUseCases) Add(user models.UserRequest) (*models.UserResponse, error) {
	if user.Name == "" {
		return nil, errors.New("username is empty")
	}

	if user.Password == "" {
		return nil, errors.New("password is empty")
	}

	if user.Email == "" {
		return nil, errors.New("email is empty")
	}

	if user.Phone == "" {
		return nil, errors.New("phone is empty")
	}

	if len(user.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	return u.r.Add(user)
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

func (u *UserUseCases) FindById(userid int64) (*models.User, error) {
	return u.r.FindById(userid)
}

func (u *UserUseCases) Update(userId int64) error {
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
