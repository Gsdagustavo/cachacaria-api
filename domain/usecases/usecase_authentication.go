package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/rules"
	"cachacariaapi/domain/status_codes"
	util2 "cachacariaapi/domain/util"
	repositories "cachacariaapi/infrastructure/datastore"
	"cachacariaapi/infrastructure/util"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AuthUseCases struct {
	repository     repositories.AuthRepository
	userRepository repositories.UserRepository
	authManager    util.AuthManager
	emailConfig    util2.EmailConfig
}

func NewAuthUseCases(repository repositories.AuthRepository, userRepository repositories.UserRepository, authManager util.AuthManager, emailConfig util2.EmailConfig) AuthUseCases {
	return AuthUseCases{
		repository:     repository,
		userRepository: userRepository,
		authManager:    authManager,
		emailConfig:    emailConfig,
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

	if !a.authManager.CheckPasswordHash(credentials.Password, user.Password) {
		return "", status_codes.LoginInvalidCredentials, nil
	}

	token, err := a.authManager.CreateToken(credentials.Email, user.ID, user.IsAdm)
	if err != nil {
		return "", status_codes.LoginFailure, errors.Join(fmt.Errorf("failed to generate auth token"), err)
	}

	return token, status_codes.LoginSuccess, nil
}

func (a AuthUseCases) RegisterUser(
	ctx context.Context,
	credentials entities.UserCredentials,
) (string, status_codes.RegisterStatusCode, error) {
	credentials.Email = util.TrimSpace(credentials.Email)
	credentials.Password = util.TrimSpace(credentials.Password)

	user, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to check user"), err)
	}

	if user != nil {
		return "", status_codes.RegisterUserAlreadyExist, nil
	}

	user, err = a.repository.GetUserByPhone(ctx, credentials.Phone)
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to check user"), err)
	}

	if user != nil {
		return "", status_codes.RegisterUserAlreadyExist, nil
	}

	if !rules.IsValidEmail(credentials.Email) {
		return "", status_codes.RegisterInvalidEmail, nil
	}

	if !rules.IsValidPassword(credentials.Password) {
		return "", status_codes.RegisterInvalidPassword, nil
	}

	credentials.Password, err = a.authManager.HashPassword(credentials.Password)
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to hash password"), err)
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to generate uuid"), err)
	}

	user = &entities.User{
		UUID:     userUUID,
		Email:    credentials.Email,
		Password: credentials.Password,
		Phone:    credentials.Phone,
		IsAdm:    credentials.IsAdm,
	}

	id, err := a.repository.AddUser(ctx, user)
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to add user"), err)
	}

	user.ID = int(id)

	token, err := a.authManager.CreateToken(credentials.Email, user.ID, user.IsAdm)
	if err != nil {
		return "", status_codes.RegisterFailure, errors.Join(fmt.Errorf("failed to generate auth token"), err)
	}

	go util2.SendAccountCreatedEmail(a.emailConfig, []string{user.Email}, *user)

	return token, status_codes.RegisterSuccess, nil
}

func (a AuthUseCases) GetUserByAuthToken(token string) (*entities.User, error) {
	payload, err := a.authManager.VerifyToken(token)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to verify token"), err)
	}

	return a.userRepository.FindById(int64(payload.UserID))
}

func (a AuthUseCases) ChangePassword(ctx context.Context, request entities.ChangePasswordRequest) (status_codes.ChangePasswordStatus, error) {
	request.CurrentPassword = util.TrimSpace(request.CurrentPassword)
	request.NewPassword = util.TrimSpace(request.NewPassword)
	request.NewPasswordConfirmation = util.TrimSpace(request.NewPasswordConfirmation)

	user, err := a.repository.GetUserByID(ctx, request.UserID)
	if err != nil {
		return status_codes.ChangePasswordError, errors.Join(fmt.Errorf("failed to check user"), err)
	}

	if user == nil {
		return status_codes.ChangePasswordInvalidUser, nil
	}

	if request.CurrentPassword == "" {
		return status_codes.ChangePasswordIncomplete, nil
	}

	if request.NewPassword == "" {
		return status_codes.ChangePasswordIncomplete, nil
	}

	if request.NewPasswordConfirmation == "" {
		return status_codes.ChangePasswordIncomplete, nil
	}

	isValidPreviousPassword := a.authManager.CheckPasswordHash(request.CurrentPassword, user.Password)
	if !isValidPreviousPassword {
		return status_codes.ChangePasswordInvalidPassword, nil
	}

	if request.NewPassword != request.NewPasswordConfirmation {
		return status_codes.ChangePasswordPasswordsDontMatch, nil
	}

	isNewPasswordEqual := a.authManager.CheckPasswordHash(request.NewPassword, user.Password)
	if isNewPasswordEqual {
		return status_codes.ChangePasswordAlreadyUsedPassword, nil
	}

	hashedNewPassword, err := a.authManager.HashPassword(request.NewPassword)
	if err != nil {
		return status_codes.ChangePasswordError, errors.Join(fmt.Errorf("failed to hash password"), err)
	}

	err = a.repository.UpdateUserPassword(ctx, int64(user.ID), hashedNewPassword)
	if err != nil {
		return status_codes.ChangePasswordError, errors.Join(fmt.Errorf("failed to update password"), err)
	}

	go util2.SendPasswordChangedEmail(a.emailConfig, []string{user.Email}, *user)

	return status_codes.ChangePasswordSuccess, nil
}
