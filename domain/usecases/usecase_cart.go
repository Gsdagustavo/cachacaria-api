package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/status_codes"
	repositories "cachacariaapi/infrastructure/datastore"
	"context"
	"errors"
	"fmt"
)

type CartUseCases struct {
	cartRepository    repositories.CartRepository
	userRepository    repositories.UserRepository
	productRepository repositories.ProductRepository
}

func NewCartUseCases(repo repositories.CartRepository, userRepository repositories.UserRepository, productRepository repositories.ProductRepository) CartUseCases {
	return CartUseCases{cartRepository: repo, productRepository: productRepository, userRepository: userRepository}
}

func (uc *CartUseCases) AddToCart(ctx context.Context, userID, productID int64, quantity int) (status_codes.AddProductItemStatus, error) {
	user, err := uc.userRepository.FindById(userID)
	if err != nil {
		return status_codes.AddProductItemStatusError, errors.Join(fmt.Errorf("failed to add product to cart"), err)
	}

	if user == nil {
		return status_codes.AddProductItemStatusInvalidUser, nil
	}

	product, err := uc.productRepository.GetProduct(productID)
	if err != nil {
		return status_codes.AddProductItemStatusError, errors.Join(fmt.Errorf("failed to add product to cart"), err)
	}

	if product == nil {
		return status_codes.AddProductItemStatusInvalidProduct, nil
	}

	if product.Stock < quantity || quantity > 999 || quantity < 1 {
		return status_codes.AddProductItemStatusInvalidQuantity, nil
	}

	err = uc.cartRepository.AddToCart(ctx, userID, productID, quantity)
	if err != nil {
		return status_codes.AddProductItemStatusError, errors.Join(fmt.Errorf("failed to add product to cart"), err)
	}

	return status_codes.AddProductItemStatusSuccess, nil
}

func (uc *CartUseCases) GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error) {

	return uc.cartRepository.GetCartItems(ctx, userID)
}

func (uc *CartUseCases) UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error {
	return uc.cartRepository.UpdateCartItem(ctx, userID, productID, quantity)
}

func (uc *CartUseCases) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	return uc.cartRepository.DeleteCartItem(ctx, userID, productID)
}

func (uc *CartUseCases) ClearCart(ctx context.Context, userID int64) error {
	return uc.cartRepository.ClearCart(ctx, userID)
}
