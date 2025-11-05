package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"context"
)

type CartUseCases struct {
	repository repositories.CartRepository
}

func NewCartUseCases(repo repositories.CartRepository) *CartUseCases {
	return &CartUseCases{repository: repo}
}

func (uc *CartUseCases) AddToCart(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repository.AddToCart(ctx, userID, productID, quantity)
}

func (uc *CartUseCases) GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error) {
	return uc.repository.GetCartItems(ctx, userID)
}

func (uc *CartUseCases) UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repository.UpdateCartItem(ctx, userID, productID, quantity)
}

func (uc *CartUseCases) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	return uc.repository.DeleteCartItem(ctx, userID, productID)
}

func (uc *CartUseCases) ClearCart(ctx context.Context, userID int64) error {
	return uc.repository.ClearCart(ctx, userID)
}
