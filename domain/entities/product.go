package entities

import (
	"mime/multipart"
	"time"
)

// Product represents a product in the catalog.
type Product struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	Price       float32  `json:"price"`
	Stock       int      `json:"stock"`
}

// CartItem links a user to a product with a quantity.
type CartItem struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Product    *Product  `json:"product,omitempty"`
	ProductID  int64     `json:"product_id"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// Order represents a completed purchase made by a user.
type Order struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	TotalAmount float64      `json:"total_amount"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	ModifiedAt  time.Time    `json:"modified_at"`
	Items       []*OrderItem `json:"items,omitempty"`
}

// OrderItem represents a single product inside an order.
type OrderItem struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	ProductID  int64     `json:"product_id"`
	Product    *Product  `json:"product,omitempty"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
type AddProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	Photos      []*multipart.FileHeader
}

type AddProductResponse struct {
	ID int64 `json:"id"`
}

type DeleteProductRequest struct {
	ID int64 `json:"id"`
}

type DeleteProductResponse struct {
	ID int64 `json:"id"`
}

type UpdateProductRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Price       *float32  `json:"price"`
	Stock       *int      `json:"stock"`
	Photos      *[]string `json:"photos"`
}

type UpdateProductResponse struct {
	ID int64 `json:"id"`
}

type PaginatedProducts struct {
	Page     int
	Offset   int
	Limit    int
	Products []Product
}
