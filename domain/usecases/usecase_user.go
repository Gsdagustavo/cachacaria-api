package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"log/slog"
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
		slog.Error("error finding user by id", "error", err.Error())
		return err
	}

	err = u.r.Delete(userId)
	if err != nil {
		slog.Error("error deleting user", "error", err.Error())
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
func (u *UserUseCases) Update(user entities.User, userId int64) error {
	if userId <= 0 {
		return nil
	}

	err := u.r.Update(user, userId)
	if err != nil {
		slog.Error("error updating user", "error", err.Error())
		return nil

	}

	return nil
}
