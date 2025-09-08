package entities

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	Type         string  `json:"type"`
	Origin       string  `json:"origin"`
	Manufacturer string  `json:"manufacturer"`
	Award        string  `json:"award"`
}

type ProductRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	Type         string  `json:"type"`
	Origin       string  `json:"origin"`
	Manufacturer string  `json:"manufacturer"`
	Award        string  `json:"award"`
}

type AddProductResponse struct {
	ID int64 `json:"id"`
}
