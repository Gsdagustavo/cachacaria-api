package repositories

import (
	"cachacariaapi/internal/models"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	Add(user models.UserRequest) (*models.UserResponse, error)
	Delete(userId int64) error
	FindByEmail(email string) (*models.User, error)
	FindById(userid int64) (*models.User, error)
	Update(user models.UserRequest, userId int64) (*models.UserResponse, error)
}
