package userusecases

import (
	"cachacariaapi/internal/http/core"
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

// GetAll users, or an error if any occurs
func (u *UserUseCases) GetAll() ([]models.User, error) {
	return u.r.GetAll()
}

// Add a user and returns a UserRespons, or an error if any occurs
func (u *UserUseCases) Add(user models.UserRequest) (*models.UserResponse, error) {

	// if any of the fields were not given, return an ApiError with code 400 and a message
	if user.Name == "" || user.Password == "" || user.Email == "" || user.Phone == "" {
		return nil, &core.ApiError{Code: 400, Message: "all fields are required"}
	}

	// if the password length is less than the required length, return an ApiError with code 400 and a message
	if len(user.Password) < 8 {
		return nil, &core.ApiError{Code: 400, Message: "password must contain at least 8 characters"}
	}

	// try to add the user calling the repository
	res, err := u.r.Add(user)

	// if any other error occurs while adding the user, return an ApiError with code 500 and a message
	if err != nil {
		return nil, &core.ApiError{Code: 500, Message: "could not add user", Err: err}
	}

	// return the response
	return res, nil
}

// Delete a user with the given userId. Return an error if any occurs
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

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*models.User, error) {
	return u.r.FindById(userid)
}

// Update TODO: add update method in the user repository
//func (u *UserUseCases) Update(userId int64) error {
//
//	_, err := u.FindById(userId)
//	if err != nil {
//		return errors.New("user not found. error: " + err.Error())
//	}
//
//	err = u.r.Delete(userId)
//	if err != nil {
//		return errors.New("user could not be deleted. error: " + err.Error())
//	}
//
//	return nil
//}
