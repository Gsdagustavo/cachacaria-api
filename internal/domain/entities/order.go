package entities

import (
	"time"
)

type Order struct {
	ID       int64     `json:"order_id"`
	UserID   string    `json:"user_id"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
	Products []Product `json:"products"`
}

type AddOrderRequest struct {
	UserID   int64     `json:"user_id"`
	Products []Product `json:"products"`
}

type AddOrderResponse struct {
	ID int64 `json:"order_id"`
}
