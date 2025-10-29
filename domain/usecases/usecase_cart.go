package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
)

type CartUseCases struct {
	repository repositories.CartRepository
}

func NewCartUseCases(repository repositories.CartRepository) *CartUseCases {
	return &CartUseCases{
		repository: repository,
	}
}

func (u *CartUseCases) GetCartByUserID(id int64) ([]entities.CartItem, error) {
	return u.repository.GetCartByUserID(id)
}
