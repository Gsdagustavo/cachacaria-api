package repositories

import (
	"cachacariaapi/internal/domain/entities"
	"mime/multipart"
)

type ProductRepository interface {
	Add(product entities.AddProductRequest, photos []*multipart.FileHeader) (*entities.AddProductResponse, error)
	GetAll() ([]entities.Product, error)
}
