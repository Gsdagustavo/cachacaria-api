package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"errors"
	"fmt"
)

type UserUseCases struct {
	r repositories.UserRepository
}

func NewUserUseCases(r repositories.UserRepository) *UserUseCases {
	return &UserUseCases{r}
}

// GetAll users, or an error if any occurs
func (u *UserUseCases) GetAll() ([]entities.User, error) {
	return u.r.GetAll()
}

// Delete a user with the given userId. Return an error if any occurs
func (u *UserUseCases) Delete(userId int64) error {
	_, err := u.FindById(userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to get user by id"), err)
	}

	err = u.r.Delete(userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to delete user"), err)
	}

	return nil
}

// FindByEmail returns the user with the given email, or an error if any occurs
func (u *UserUseCases) FindByEmail(email string) (*entities.User, error) {
	user, err := u.r.FindByEmail(email)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to find user by email"), err)
	}

	return user, nil
}

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*entities.User, error) {
	user, err := u.r.FindById(userid)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to find user by id"), err)
	}

	return user, nil
}

// Update a user from the database from the given UserRequest and userId
// Returns a UserResponse oran error if any occurs
func (u *UserUseCases) Update(user entities.User, userId int64) error {
	if userId <= 0 {
		return nil
	}

	err := u.r.Update(user, userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to update user"), err)
	}

	return nil
}
