package userusecases

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/core"
	"errors"
	"log"
	"regexp"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
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
	log.Printf("add user called on user usecases. request: %v", req)
	log.Printf("password: %v", req.Password)

	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	user, err := u.r.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, core.ErrConflict
	}

	if err := validatePhone(req.Phone); err != nil {
		return nil, err
	}

	// hash password

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

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		log.Printf("invalid email: %s", email)
		return core.ErrInvalidEmail
	}

	return nil
}

func validatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		log.Printf("invalid password: %s", password)
		return core.ErrInvalidPassword
	}

	return nil
}

func validatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		log.Printf("invalid phone: %s", phone)
		return core.ErrInvalidPhoneNumber
	}

	return nil
}
