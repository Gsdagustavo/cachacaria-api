package repositories

import (
	"cachacariaapi/domain/entities"
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	AddUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*entities.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUserPassword(ctx context.Context, userID int64, newPassword string) error
}

type UserRepository interface {
	GetAll() ([]entities.User, error)
	Add(user entities.User) error
	Delete(userId int64) error
	FindByEmail(email string) (*entities.User, error)
	FindById(userid int64) (*entities.User, error)
	Update(user entities.User) error
}

type ProductRepository interface {
	AddProduct(product entities.AddProductRequest) (int64, error)
	AddProductPhotos(productID int64, filenames []string) error
	GetAll(limit, offset int) ([]entities.Product, error)
	GetProduct(id int64) (*entities.Product, error)
	DeleteProduct(id int64) error
	UpdateProduct(id int64, product entities.UpdateProductRequest) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, userID, productID int64, quantity int) error
	GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error)
	UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error
	DeleteCartItem(ctx context.Context, userID, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}
