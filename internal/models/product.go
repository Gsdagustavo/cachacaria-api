package models

import "github.com/google/uuid"

type Product struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float32   `json:"price"`
	Type         string    `json:"type"`
	Origin       string    `json:"origin"`
	Manufacturer string    `json:"manufacturer"`
	Award        string    `json:"award"`
}

type AddProductRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	Type         string  `json:"type"`
	Origin       string  `json:"origin"`
	Manufacturer string  `json:"manufacturer"`
	Award        string  `json:"award"`
}

type AddProductResponse struct {
	ID uuid.UUID `json:"id"`
}
