package usecases

import "cachacariaapi/internal/models"

type UserUseCases interface {
	GetAll() []models.User
	Add(user models.AddUserRequest) (*models.AddUserResponse, error)
	Delete(userId int64) error
	FindById(userid int64) (*models.User, error)
}
