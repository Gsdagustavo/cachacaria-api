package entities

type Product struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Reviews     []Review `json:"reviews"`
	Price       float32  `json:"price"`
	Stock       int      `json:"stock"`
}

type AddProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Price float32 `json:"price"`
	Stock int     `json:"stock"`
}

type AddProductResponse struct {
	ID int64 `json:"id"`
}
