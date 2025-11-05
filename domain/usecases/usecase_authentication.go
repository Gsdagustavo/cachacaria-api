package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"cachacariaapi/domain/rules"
	"cachacariaapi/domain/status_codes"
	"cachacariaapi/infrastructure/util"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AuthUseCases struct {
	repository     repositories.AuthRepository
	userRepository repositories.UserRepository
	crypt          util.Crypt
}

func NewAuthUseCases(repository repositories.AuthRepository, userRepository repositories.UserRepository, crypt util.Crypt) *AuthUseCases {
	return &AuthUseCases{
		repository:     repository,
		userRepository: userRepository,
		crypt:          crypt,
	}
}

func (a AuthUseCases) AttemptLogin(
	ctx context.Context,
	credentials entities.UserCredentials,
) (string, status_codes.LoginStatusCode, error) {
	user, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return "", status_codes.LoginFailure, errors.Join(fmt.Errorf("failed to get user by email"), err)
	}

	if user == nil {
		return "", status_codes.LoginUserNotFound, nil
	}

	if !a.crypt.CheckPasswordHash(credentials.Password, user.Password) {
		return "", status_codes.LoginInvalidCredentials, nil
	}

	token, err := a.crypt.GenerateAuthToken(credentials.Email, user.ID, user.IsAdm)
	if err != nil {
		return "", status_codes.LoginFailure, errors.Join(fmt.Errorf("failed to generate auth token"), err)
	}

	return token, status_codes.LoginSuccess, nil
}

func (a AuthUseCases) RegisterUser(
	ctx context.Context,
	credentials entities.UserCredentials,
) (status_codes.RegisterStatusCode, error) {
	user, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to check user"), err)
	}

	if user != nil {
		return status_codes.RegisterUserAlreadyExist, nil
	}

	credentials.Email = util.TrimSpace(credentials.Email)
	credentials.Password = util.TrimSpace(credentials.Password)

	if !rules.IsValidEmail(credentials.Email) {
		return status_codes.RegisterInvalidEmail, nil
	}

	if !rules.IsValidPassword(credentials.Password) {
		return status_codes.RegisterInvalidPassword, nil
	}

	credentials.Password, err = a.crypt.HashPassword(credentials.Password)
	if err != nil {
		return status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to hash password"), err)
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to generate uuid"), err)
	}

	user = &entities.User{
		UUID:     userUUID,
		Email:    credentials.Email,
		Password: credentials.Password,
		Phone:    credentials.Phone,
		IsAdm:    credentials.IsAdm,
	}

	err = a.repository.AddUser(ctx, user)
	if err != nil {
		return status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to add user"), err)
	}

	return status_codes.RegisterSuccess, nil
}
func (a AuthUseCases) GetUserByAuthToken(token string) (*entities.User, error) {
	payload, err := a.crypt.VerifyAuthToken(token)
	if err != nil {
		return nil, err
	}

	return a.userRepository.FindById(int64(payload.UserID))
}
