package repositories

import "cachacariaapi/domain/entities"

type AddProductToCartRequest struct {
	product  entities.Product
	user     entities.User
	quantity int
}

type DeleteProductFromCartRequest struct {
	product  entities.Product
	user     entities.User
	quantity int
}

type UpdateCartRequest struct {
	product  entities.Product
	user     entities.User
	quantity int
}

type CartRepository interface {
	GetCartByUserID(id int64) ([]entities.CartItem, error)
	AddProductToCart(userID, productID int64, quantity int) error
	RemoveProductFromCart(userID, productID int64) error
}
