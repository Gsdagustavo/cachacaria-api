package repositories

import (
	"cachacariaapi/domain/entities"
	repositories "cachacariaapi/infrastructure/datastore"
	"context"
	"database/sql"
)

type MysqlOrderRepository struct {
	DB *sql.DB
}

func NewMysqlOrderRepository(db *sql.DB) repositories.OrderRepository {
	return &MysqlOrderRepository{
		DB: db,
	}
}

func (r *MysqlOrderRepository) GetOrders(ctx context.Context, userID int64) ([]entities.Order, error) {
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
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	var currentOrder *entities.Order

	for rows.Next() {
		var order entities.Order
		var item entities.OrderItem

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CreatedAt,
			&order.ModifiedAt,
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		// Check if we need to add a new order to the list
		if currentOrder == nil || currentOrder.ID != order.ID {
			if currentOrder != nil {
				orders = append(orders, *currentOrder)
			}
			currentOrder = &order
		}

		// Add order item to the current order
		currentOrder.Items = append(currentOrder.Items, item)
	}

	// Append the last order
	if currentOrder != nil {
		orders = append(orders, *currentOrder)
	}

	return orders, nil
}

func (r *MysqlOrderRepository) AddOrder(ctx context.Context, order entities.Order) error {
	// Start a transaction to ensure both the order and order items are added atomically
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert the order into the 'orders' table
	const orderQuery = `
		INSERT INTO orders (user_id) 
		VALUES (?)`
	result, err := tx.ExecContext(ctx, orderQuery, order.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the generated order ID
	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert each order item into the 'order_items' table
	const orderItemQuery = `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES (?, ?, ?, ?)`

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, orderItemQuery, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}
