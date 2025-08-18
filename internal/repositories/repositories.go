package repositories

import (
	"cachacariaapi/internal/models"
)

type UserRepository interface {
	GetAll() []models.User
	Add(user models.AddUserRequest) (*models.AddUserResponse, error)
}

type ProductRepository interface {
	GetAll() []models.Product
	Add(user models.AddProductRequest) (models.AddProductResponse, error)
}
