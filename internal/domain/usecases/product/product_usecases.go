package product

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/infrastructure/persistence"
	"cachacariaapi/internal/interfaces/http/core"
	"errors"
	"mime/multipart"
	"strings"
)

type ProductUseCases struct {
	r *persistence.MySQLProductRepository
}

func NewProductUseCases(r *persistence.MySQLProductRepository) *ProductUseCases {
	return &ProductUseCases{r}
}

func (u *ProductUseCases) Add(req entities.AddProductRequest, photos []*multipart.FileHeader) (*entities.AddProductResponse, error) {
	if req.Price <= 0 || len(strings.Trim(req.Name, " ")) == 0 || req.Stock < 0 {
		return nil, core.ErrBadRequest
	}

	res, err := u.r.Add(req, photos)
	if err != nil {
		if errors.Is(err, core.ErrBadRequest) {
			return nil, core.ErrBadRequest
		}

		if errors.Is(err, core.ErrConflict) {
			return nil, core.ErrConflict
		}

		return nil, core.ErrInternal
	}

	return res, nil
}

func (u *ProductUseCases) GetAll() ([]entities.Product, error) {
	return u.r.GetAll()
}
