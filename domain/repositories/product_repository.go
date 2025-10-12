package repositories

import (
	"cachacariaapi/domain/entities"
)

type ProductRepository interface {
	AddProduct(product entities.AddProductRequest) (*entities.AddProductResponse, error)
	AddProductPhotos(productID int64, filenames []string) error
	GetAll(limit, offset int) ([]entities.Product, error)
	GetProduct(id int64) (*entities.Product, error)
	DeleteProduct(id int64) (*entities.DeleteProductResponse, error)

	UpdateProduct(
		id int64,
		product entities.UpdateProductRequest,
	) (*entities.UpdateProductResponse, error)
}
