package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/status_codes"
	"cachacariaapi/domain/util"
	repositories "cachacariaapi/infrastructure/datastore"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProductUseCases struct {
	productRepository repositories.ProductRepository
	baseURL           string
}

func NewProductUseCases(productRepository repositories.ProductRepository, baseURL string) ProductUseCases {
	return ProductUseCases{productRepository, baseURL}
}

func (u *ProductUseCases) AddProduct(
	req entities.AddProductRequest,
) (status_codes.AddProductStatus, error) {
	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		return status_codes.AddProductStatusInvalidName, nil
	}

	if req.Price <= 0 || req.Stock < 0 {
		return status_codes.AddProductStatusInvalidPrice, nil
	}

	if req.Stock < 0 {
		return status_codes.AddProductStatusInvalidStock, nil
	}

	for _, fileHeader := range req.Photos {
		if err := validateImageType(fileHeader); err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to validate image type"), err)
		}
	}

	id, err := u.productRepository.AddProduct(req)
	if err != nil {
		return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to add product"), err)
	}

	var filenames []string
	for _, fileHeader := range req.Photos {
		if err = validateImageType(fileHeader); err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to validate image type"), err)
		}

		src, err := fileHeader.Open()
		if err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to open image file"), err)
		}

		filename := fmt.Sprintf(
			"product_%d_%d%s",
			id,
			time.Now().UnixNano(),
			filepath.Ext(fileHeader.Filename),
		)
		filePath := filepath.Join("/app/images", filename)

		dst, err := os.Create(filePath)
		if err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to create file"), err)
		}

		_, err = io.Copy(dst, src)
		src.Close()
		if err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("failed to copy file"), err)
		}

		filenames = append(filenames, filename)
	}

	if len(filenames) > 0 {
		if err = u.productRepository.AddProductPhotos(id, filenames); err != nil {
			return status_codes.AddProductStatusError, errors.Join(fmt.Errorf("faield to add product photos"), err)
		}
	}

	return status_codes.AddProductStatusSuccess, nil
}

func (u *ProductUseCases) GetAll(limit, page int) ([]entities.Product, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	products, err := u.productRepository.GetAll(limit, offset)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to get products"), err)
	}

	for _, product := range products {
		for i := range product.Photos {
			product.Photos[i] = util.GetProductImageURL(product.Photos[i], u.baseURL)
		}
	}

	return products, nil
}

func (u *ProductUseCases) GetProduct(id int64) (*entities.Product, error) {
	product, err := u.productRepository.GetProduct(id)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to get product"), err)
	}

	for i := range product.Photos {
		product.Photos[i] = util.GetProductImageURL(product.Photos[i], u.baseURL)
	}

	return product, nil
}

func (u *ProductUseCases) DeleteProduct(id int64) (status_codes.DeleteProductStatus, error) {
	prod, err := u.productRepository.GetProduct(id)
	if err != nil {
		return status_codes.DeleteProductStatusError, errors.Join(fmt.Errorf("failed to get product"), err)
	}

	if prod == nil {
		return status_codes.DeleteProductStatusNotFound, nil
	}

	for _, filename := range prod.Photos {
		filePath := filepath.Join("/app/images", filename)

		if err = os.Remove(filePath); err != nil {
			return status_codes.DeleteProductStatusError, entities.ErrNotFound
		}
	}

	err = u.productRepository.DeleteProduct(id)
	if err != nil {
		return status_codes.DeleteProductStatusError, errors.Join(fmt.Errorf("failed to delete product"), err)
	}

	return status_codes.DeleteProductStatusSuccess, nil
}

func (u *ProductUseCases) UpdateProduct(
	id int64,
	product entities.UpdateProductRequest,
) (status_codes.UpdateProductStatus, error) {
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

	prod, err := u.GetProduct(id)
	if err != nil {
		return status_codes.UpdateProductStatusError, errors.Join(fmt.Errorf("failed to get product"), err)
	}

	if prod == nil {
		return status_codes.UpdateProductStatusNotFound, nil
	}

	if product.Name == "" {
		return status_codes.UpdateProductStatusInvalidName, nil
	}

	if product.Description == "" {
		return status_codes.UpdateProductStatusInvalidDescription, nil
	}

	if product.Stock < 0 || product.Stock > 999 {
		return status_codes.UpdateProductStatusInvalidStock, nil
	}

	if product.Price < 0 || product.Stock > 999 {
		return status_codes.UpdateProductStatusInvalidStock, nil
	}

	err = u.productRepository.UpdateProduct(id, product)
	if err != nil {
		return status_codes.UpdateProductStatusError, errors.Join(fmt.Errorf("failed to update product"), err)
	}

	return status_codes.UpdateProductStatusSuccess, nil
}

func validateImageType(header *multipart.FileHeader) error {
	src, err := header.Open()
	if err != nil {
		return errors.Join(fmt.Errorf("failed to open image file"), err)
	}

	defer src.Close()

	buf := make([]byte, 512)
	if _, err := src.Read(buf); err != nil {
		return errors.Join(fmt.Errorf("failed to read image bytes"), err)
	}

	contentType := http.DetectContentType(buf)

	if contentType != "image/jpeg" && contentType != "image/png" {
		return entities.ErrBadRequest
	}

	return nil
}
