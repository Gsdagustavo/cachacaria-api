package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID       uuid.UUID `json:"order_id"`
	UserID   string    `json:"user_id"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
	Products []Product `json:"products"`
}

type AddOrderRequest struct {
	UserID   string    `json:"user_id"`
	Products []Product `json:"products"`
}

type AddOrderResponse struct {
	ID uuid.UUID `json:"order_id"`
}
