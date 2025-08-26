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
func (u *UserUseCases) Add(req models.RegisterRequest) (*models.UserResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	res, err := u.r.Add(req)
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

func (u *UserUseCases) FindByEmail(email string) (*models.User, error) {
	return u.r.FindByEmail(email)
}

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*models.User, error) {
	return u.r.FindById(userid)
}

// Update a user from the database from the given UserRequest and userId
// Returns a UserResponse oran error if any occurs
func (u *UserUseCases) Update(user models.UserRequest, userId int64) (*models.UserResponse, error) {
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
	if req.Password == "" || req.Email == "" || req.Phone == "" {
		return core.ErrBadRequest
	}

	if len(req.Password) < 8 {
		return core.ErrBadRequest
	}

	return nil
}

func validateRegisterRequest(req models.RegisterRequest) error {
	if req.Password == "" || req.Email == "" || req.Phone == "" {
		return core.ErrBadRequest
	}

	return nil
}
