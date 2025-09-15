package product

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/core"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProductUseCases struct {
	r *persistence.MySQLProductRepository
}

func NewProductUseCases(r *persistence.MySQLProductRepository) *ProductUseCases {
	return &ProductUseCases{r}
}

func (u *ProductUseCases) AddProduct(req entities.AddProductRequest, uploadedFiles []*multipart.FileHeader) (*entities.AddProductResponse, error) {
	if req.Price <= 0 || len(strings.TrimSpace(req.Name)) == 0 || req.Stock < 0 {
		return nil, core.ErrBadRequest
	}

	res, err := u.r.AddProduct(req)
	if err != nil {
		return nil, err
	}
	productID := res.ID

	var filenames []string
	for _, fileHeader := range uploadedFiles {
		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		filename := fmt.Sprintf("product_%d_%d%s", productID, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		filePath := filepath.Join("/app/images", filename)

		dst, err := os.Create(filePath)
		if err != nil {
			src.Close()
			return nil, err
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()
		if err != nil {
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

func (u *ProductUseCases) GetAll() ([]entities.Product, error) {
	return u.r.GetAll()
}

func (u *ProductUseCases) GetProduct(id int64) (*entities.Product, error) {
	return u.r.GetProduct(id)
}
