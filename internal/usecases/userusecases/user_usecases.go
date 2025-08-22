package userusecases

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/repositories/userrepository"
	"errors"
	"log"
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
	if err := validateUserRequest(user); err != nil {
		return nil, err
	}

	if len(user.Password) < 8 {
		return nil, core.ErrBadRequest
	}

	res, err := u.r.Add(user)
	if err != nil {
		if errors.Is(err, core.ErrConflict) {
			return nil, core.ErrConflict
		}

		return nil, core.ErrInternal
	}

	return res, nil
}

// Delete a user with the given userId. Return an error if any occurs
func (u *UserUseCases) Delete(userId int64) error {
	_, err := u.FindById(userId)
	if err != nil {
		return err
	}

	err = u.r.Delete(userId)
	if err != nil {
		return err
	}

	return nil
}

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*models.User, error) {
	return u.r.FindById(userid)
}

// Update TODO: add update method in the user repository
func (u *UserUseCases) Update(user models.UserRequest, userId int64) (*models.UserResponse, error) {
	if err := validateUserRequest(user); err != nil {
		return nil, err
	}

	if userId <= 0 {
		return nil, core.ErrBadRequest
	}

	res, err := u.r.Update(user, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateUserRequest(req models.UserRequest) error {
	if req.Name == "" || req.Password == "" || req.Email == "" || req.Phone == "" {
		log.Printf("returning bad request from usecases")
		return core.ErrBadRequest
	}

	return nil
}
