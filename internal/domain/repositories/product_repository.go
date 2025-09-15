package repositories

import (
	"cachacariaapi/internal/domain/entities"
	"mime/multipart"
)

type ProductRepository interface {
	AddProduct(product entities.AddProductRequest) (*entities.AddProductResponse, error)
	AddProductPhotos(photos []*multipart.FileHeader) error
	GetAll() ([]entities.Product, error)
	GetProduct(id int64) (*entities.Product, error)
}
