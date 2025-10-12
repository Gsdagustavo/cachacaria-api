package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProductUseCases struct {
	r       repositories.ProductRepository
	baseURL string
}

func NewProductUseCases(r repositories.ProductRepository, baseURL string) *ProductUseCases {
	return &ProductUseCases{r, baseURL}
}

func (u *ProductUseCases) AddProduct(
	req entities.AddProductRequest,
) (*entities.AddProductResponse, error) {
	if req.Price <= 0 || len(strings.TrimSpace(req.Name)) == 0 || req.Stock < 0 {
		return nil, entities.ErrBadRequest
	}

	for _, fileHeader := range req.Photos {
		if err := validateImageType(fileHeader); err != nil {
			slog.Error("error validating image type", "error", err.Error())
			return nil, err
		}
	}

	res, err := u.r.AddProduct(req)
	if err != nil {
		slog.Error("error adding product", "error", err.Error())
		return nil, err
	}
	productID := res.ID

	var filenames []string
	for _, fileHeader := range req.Photos {
		if err := validateImageType(fileHeader); err != nil {
			slog.Error("error validating image type", "error", err.Error())
		}

		src, err := fileHeader.Open()
		if err != nil {
			slog.Error("error opening file header", "error", err.Error())
		}

		filename := fmt.Sprintf(
			"product_%d_%d%s",
			productID,
			time.Now().UnixNano(),
			filepath.Ext(fileHeader.Filename),
		)
		filePath := filepath.Join("/app/images", filename)

		dst, err := os.Create(filePath)
		if err != nil {
			slog.Error("error creating", "error", err.Error())
			return nil, err
		}

		_, err = io.Copy(dst, src)
		src.Close()
		if err != nil {
			slog.Error("error copying file", "error", err.Error())
			return nil, err
		}

		filenames = append(filenames, filename)
	}

	if len(filenames) > 0 {
		if err = u.r.AddProductPhotos(productID, filenames); err != nil {
			return nil, err
		}
	}

	return &entities.AddProductResponse{ID: productID}, nil
}

func (u *ProductUseCases) GetAll(limit, page int) ([]entities.Product, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	products, err := u.r.GetAll(limit, offset)
	if err != nil {
		slog.Error("error getting all products", "error", err.Error())
		return nil, err
	}

	for _, product := range products {
		for i := range product.Photos {
			product.Photos[i] = u.buildProductImageURL(product.Photos[i])
		}
	}

	return products, nil
}

func (u *ProductUseCases) GetProduct(id int64) (*entities.Product, error) {
	product, err := u.r.GetProduct(id)
	if err != nil {
		slog.Error("error getting product", "error", err.Error())
		return nil, err
	}

	for i := range product.Photos {
		product.Photos[i] = u.buildProductImageURL(product.Photos[i])
	}

	return product, nil
}

func (u *ProductUseCases) DeleteProduct(id int64) (*entities.DeleteProductResponse, error) {
	prod, err := u.r.GetProduct(id)
	if err != nil {
		slog.Error("error getting product for deletion check", "error", err.Error())
		return nil, err
	}

	if prod == nil {
		return nil, entities.ErrNotFound
	}

	for _, filename := range prod.Photos {
		filePath := filepath.Join("/app/images", filename)

		if err := os.Remove(filePath); err != nil {
			slog.Error(
				"error deleting product photo file",
				"filePath",
				filePath,
				"error",
				err.Error(),
			)
		}
	}

	res, err := u.r.DeleteProduct(id)
	if err != nil {
		slog.Error("error deleting product from database", "error", err.Error())
		return nil, err
	}

	return res, nil
}

func (u *ProductUseCases) UpdateProduct(
	id int64,
	product entities.UpdateProductRequest,
) (*entities.UpdateProductResponse, error) {
	prod, err := u.GetProduct(id)
	if err != nil {
		slog.Error("error getting product", "error", err.Error())
		return nil, err
	}

	if prod == nil {
		return nil, entities.ErrNotFound
	}

	res, err := u.r.UpdateProduct(id, product)
	if err != nil {
		slog.Error("error updating product", "error", err.Error())
		return nil, err
	}

	return res, nil
}

func validateImageType(header *multipart.FileHeader) error {
	src, err := header.Open()
	if err != nil {
		slog.Error("error opening file", "error", err.Error())
		return err
	}

	defer src.Close()

	buf := make([]byte, 512)
	if _, err := src.Read(buf); err != nil {
		slog.Error("error reading image bytes", "error", err.Error())
	}

	contentType := http.DetectContentType(buf)

	if contentType != "image/jpeg" && contentType != "image/png" {
		return entities.ErrBadRequest
	}

	return nil
}

func (u *ProductUseCases) buildProductImageURL(filename string) string {
	return fmt.Sprintf("%s/images/%s", u.baseURL, filename)
}
