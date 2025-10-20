package repositories

import (
	"cachacariaapi/domain/entities"
)

type ProductRepository interface {
	AddProduct(product entities.AddProductRequest) (int64, error)
	AddProductPhotos(productID int64, filenames []string) error
	GetAll(limit, offset int) ([]entities.Product, error)
	GetProduct(id int64) (*entities.Product, error)
	DeleteProduct(id int64) error

	UpdateProduct(
		id int64,
		product entities.UpdateProductRequest,
	) error
}
