package entities

import "mime/multipart"

type Product struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	Reviews     []Review `json:"reviews"`
	Price       float32  `json:"price"`
	Stock       int      `json:"stock"`
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
