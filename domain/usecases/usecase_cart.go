package usecases

import (
	"cachacariaapi/domain/entities"
	"context"
)

type CartUseCases struct {
	repo CartRepository
}

type CartRepository interface {
	AddToCart(ctx context.Context, userID, productID int64, quantity int) error
	GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error)
	UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error
	DeleteCartItem(ctx context.Context, userID, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}

func NewCartUseCases(repo CartRepository) *CartUseCases {
	return &CartUseCases{repo: repo}
}

func (uc *CartUseCases) AddToCart(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repo.AddToCart(ctx, userID, productID, quantity)
}

func (uc *CartUseCases) GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error) {
	return uc.repo.GetCartItems(ctx, userID)
}

func (uc *CartUseCases) UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.repo.UpdateCartItem(ctx, userID, productID, quantity)
}

func (uc *CartUseCases) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	return uc.repo.DeleteCartItem(ctx, userID, productID)
}

func (uc *CartUseCases) ClearCart(ctx context.Context, userID int64) error {
	return uc.repo.ClearCart(ctx, userID)
}
