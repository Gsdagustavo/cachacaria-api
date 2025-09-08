package entities

import "image"

type Product struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Photos      []image.Image `json:"photos"`
	Reviews     []Review      `json:"reviews"`
	Description string        `json:"description"`
	Price       float32       `json:"price"`
	Stock       int           `json:"stock"`
}

type AddProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}

type AddProductResponse struct {
	ID int64 `json:"id"`
}
