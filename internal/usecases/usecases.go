package usecases

import "cachacariaapi/internal/models"

type UserUseCases interface {
	GetAll() ([]models.User, error)
	Add(user models.RegisterRequest) (*models.UserResponse, error)
	Delete(userId int64) error
	FindByEmail(email string) (*models.User, error)
	FindById(userid int64) (*models.User, error)
	Update(user models.UserRequest, userId int64) (*models.UserResponse, error)
}
