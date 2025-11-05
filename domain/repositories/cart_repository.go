package repositories

import (
	"cachacariaapi/domain/entities"
	"context"
)

type CartRepository interface {
	AddToCart(ctx context.Context, userID, productID int64, quantity int) error
	GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error)
	UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error
	DeleteCartItem(ctx context.Context, userID, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}
