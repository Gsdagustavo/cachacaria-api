package userusecases

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/core"
	"errors"
)

type UserUseCases struct {
	r *persistence.MySQLUserRepository
}

func NewUserUseCases(r *persistence.MySQLUserRepository) *UserUseCases {
	return &UserUseCases{r}
}

// GetAll users, or an error if any occurs
func (u *UserUseCases) GetAll() ([]entities.User, error) {
	return u.r.GetAll()
}

// Add a user and returns a UserResponse, or an error if any occurs
func (u *UserUseCases) Add(req entities.RegisterRequest) (*entities.UserResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	user, err := u.r.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, core.ErrConflict
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

// FindByEmail returns the user with the given email, or an error if any occurs
func (u *UserUseCases) FindByEmail(email string) (*entities.User, error) {
	return u.r.FindByEmail(email)
}

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*entities.User, error) {
	return u.r.FindById(userid)
}

// Update a user from the database from the given UserRequest and userId
// Returns a UserResponse oran error if any occurs
func (u *UserUseCases) Update(user entities.UserRequest, userId int64) (*entities.UserResponse, error) {
	if userId <= 0 {
		return nil, core.ErrBadRequest
	}

	res, err := u.r.Update(user, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateRegisterRequest(req entities.RegisterRequest) error {
	if req.Password == "" || req.Email == "" || req.Phone == "" {
		return core.ErrBadRequest
	}

	return nil
}
