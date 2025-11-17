package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/rules"
	"cachacariaapi/domain/status_codes"
	repositories "cachacariaapi/infrastructure/datastore"
	"cachacariaapi/infrastructure/util"
	"errors"
	"fmt"
)

type UserUseCases struct {
	userRepository repositories.UserRepository
	authRepository repositories.AuthRepository
	authManager    util.AuthManager
}

func NewUserUseCases(userRepository repositories.UserRepository, authRepository repositories.AuthRepository, authManager util.AuthManager) UserUseCases {
	return UserUseCases{
		authRepository: authRepository,
		userRepository: userRepository,
		authManager:    authManager,
	}
}

// GetAll users, or an error if any occurs
func (u *UserUseCases) GetAll() ([]entities.User, error) {
	return u.userRepository.GetAll()
}

// Delete a user with the given userId. Return an error if any occurs
func (u *UserUseCases) Delete(userId int64) error {
	_, err := u.FindById(userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to get user by id"), err)
	}

	err = u.userRepository.Delete(userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to delete user"), err)
	}

	return nil
}

// FindByEmail returns the user with the given email, or an error if any occurs
func (u *UserUseCases) FindByEmail(email string) (*entities.User, error) {
	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to find user by email"), err)
	}

	return user, nil
}

// FindById returns the user with the given userId, or an error if any occurs
func (u *UserUseCases) FindById(userid int64) (*entities.User, error) {
	user, err := u.userRepository.FindById(userid)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to find user by id"), err)
	}

	return user, nil
}

// Update a user from the database from the given UserRequest and userId
// Returns a UserResponse oran error if any occurs
func (u *UserUseCases) Update(user entities.User) (status_codes.UpdateUserStatus, error) {
	user.Email = util.TrimSpace(user.Email)
	user.Phone = util.TrimSpace(user.Phone)

	existingUser, err := u.userRepository.FindByEmail(user.Email)
	if err != nil {
		return status_codes.UpdateUserFailure, errors.Join(fmt.Errorf("failed to check user"), err)
	}

	if existingUser != nil {
		return status_codes.UpdateUserEmailAlreadyExists, nil
	}

	if !rules.IsValidEmail(user.Email) {
		return status_codes.UpdateUserInvalidEmail, nil
	}

	err = u.userRepository.Update(user)
	if err != nil {
		return status_codes.UpdateUserFailure, errors.Join(fmt.Errorf("failed to update user"), err)
	}

	return status_codes.UpdateUserSuccess, nil
}
