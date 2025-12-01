package repositories

import (
	"cachacariaapi/domain/entities"
	repositories "cachacariaapi/infrastructure/datastore"
	"context"
	"database/sql"
)

type MYSQLOrderRepository struct {
	DB *sql.DB
}

func NewMYSQLOrderRepository(db *sql.DB) repositories.OrderRepository {
	return &MYSQLOrderRepository{
		DB: db,
	}
}

func (r *MYSQLOrderRepository) CreateOrder(ctx context.Context, userID int64) (int64, error) {
	const query = `
        INSERT INTO orders (user_id)
        VALUES (?)
    `

	result, err := r.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *MYSQLOrderRepository) AddOrderItem(ctx context.Context, orderID, productID int64, quantity int) error {
	const query = `
        INSERT INTO order_items (order_id, product_id, quantity, price)
        SELECT ?, id, ?, price
        FROM products
        WHERE id = ?
    `
	_, err := r.DB.ExecContext(ctx, query, orderID, quantity, productID)
	return err
}

func (r *MYSQLOrderRepository) GetOrders(ctx context.Context, userID int64) ([]entities.Order, error) {
	const query = `
		SELECT 
			o.id,
			o.user_id,
			o.created_at,
			o.modified_at,
			oi.id AS item_id,
			oi.product_id,
			oi.quantity,
			oi.price
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = ?
		ORDER BY o.id
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	var currentOrder *entities.Order

	for rows.Next() {
		var o entities.Order
		var item entities.OrderItem

		err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.CreatedAt,
			&o.ModifiedAt,
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		if currentOrder == nil || currentOrder.ID != o.ID {
			if currentOrder != nil {
				orders = append(orders, *currentOrder)
			}
			currentOrder = &o
		}

		if item.ID != 0 {
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}

	if currentOrder != nil {
		orders = append(orders, *currentOrder)
	}

	return orders, nil
}
