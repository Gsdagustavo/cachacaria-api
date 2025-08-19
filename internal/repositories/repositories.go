package repositories

import (
	"cachacariaapi/internal/models"
)

type UserRepository interface {
	GetAll() []models.User
	Add(user models.AddUserRequest) (*models.AddUserResponse, error)
	Delete(userId int64) error
	FindById(userid int64) (*models.User, error)
}

type ProductRepository interface {
	GetAll() []models.Product
	Add(user models.AddProductRequest) (models.AddProductResponse, error)
}
