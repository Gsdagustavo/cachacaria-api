package entities

import "time"

type CartItem struct {
	ID         int64
	UserID     int64
	ProductID  int64
	Product    *Product
	Quantity   int
	CreatedAt  time.Time
	ModifiedAt time.Time
}
