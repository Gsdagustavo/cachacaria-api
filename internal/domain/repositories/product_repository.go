package repositories

import "cachacariaapi/internal/domain/entities"

type ProductRepository interface {
	GetAll() ([]entities.Product, error)
	Add(product entities.AddProductRequest) (*entities.AddProductResponse, error)
}
