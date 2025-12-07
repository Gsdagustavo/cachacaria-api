package usecases

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/status_codes"
	"cachacariaapi/domain/util"
	repositories "cachacariaapi/infrastructure/datastore"
	"context"
	"errors"
	"fmt"
)

type CartUseCases struct {
	cartRepository    repositories.CartRepository
	userRepository    repositories.UserRepository
	productRepository repositories.ProductRepository
	orderRepository   repositories.OrderRepository
	baseURL           string
}

func NewCartUseCases(
	repo repositories.CartRepository,
	userRepository repositories.UserRepository,
	productRepository repositories.ProductRepository,
	orderRepository repositories.OrderRepository,
	baseURL string,
) CartUseCases {
	return CartUseCases{
		cartRepository:    repo,
		productRepository: productRepository,
		userRepository:    userRepository,
		orderRepository:   orderRepository,
		baseURL:           baseURL,
	}
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
	items, err := uc.cartRepository.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Product != nil && len(item.Product.Photos) > 0 {
			photos := make([]string, len(item.Product.Photos))
			for i, filename := range item.Product.Photos {
				photos[i] = util.GetProductImageURL(filename, uc.baseURL)
			}
			item.Product.Photos = photos
		}
	}

	return items, nil
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

func (uc *CartUseCases) BuyItems(ctx context.Context, userID int64) (status_codes.BuyProductsStatus, error) {
	items, err := uc.cartRepository.GetCartItems(ctx, userID)
	if err != nil {
		return status_codes.BuyProductsStatusError, err
	}

	if len(items) == 0 {
		return status_codes.BuyProductsStatusCartEmpty, nil
	}

	for _, item := range items {
		product, err := uc.productRepository.GetProduct(item.ProductID)
		if err != nil {
			return status_codes.BuyProductsStatusError, err
		}

		if product == nil {
			return status_codes.BuyProductsStatusInvalidProduct, nil
		}

		if product.Stock < item.Quantity {
			return status_codes.BuyProductsStatusOutOfStock, nil
		}
	}

	err = uc.cartRepository.ClearCart(ctx, userID)
	if err != nil {
		return status_codes.BuyProductsStatusError, err
	}

	return status_codes.BuyProductsStatusSuccess, nil
}

func (uc *CartUseCases) GetOrders(ctx context.Context, userID int64) ([]entities.Order, error) {
	user, err := uc.userRepository.FindById(userID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to find user"), err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	orders, err := uc.orderRepository.GetOrders(ctx, userID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to fetch orders"), err)
	}

	for i := range orders {
		for j := range orders[i].Items {
			item := &orders[i].Items[j]

			product, err := uc.productRepository.GetProduct(item.ProductID)
			if err != nil {
				return nil, err
			}

			item.Product = product

			if product != nil && len(product.Photos) > 0 {
				photos := make([]string, len(product.Photos))
				for k, filename := range product.Photos {
					photos[k] = util.GetProductImageURL(filename, uc.baseURL)
				}
				item.Product.Photos = photos
			}
		}
	}

	return orders, nil
}
