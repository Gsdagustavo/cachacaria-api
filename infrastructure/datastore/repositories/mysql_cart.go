package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/repositories"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type MySQLCartRepository struct {
	DB *sql.DB
}

func NewMySQLCartRepository(db *sql.DB) repositories.CartRepository {
	return &MySQLCartRepository{DB: db}
}

func (repo *MySQLCartRepository) AddToCart(ctx context.Context, userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRowContext(ctx, `
        SELECT quantity FROM carts_products
        WHERE user_id = ? AND product_id = ? FOR UPDATE;
    `, userID, productID).Scan(&currentQty)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = tx.ExecContext(ctx, `
            INSERT INTO carts_products (user_id, product_id, quantity, created_at, modified_at)
            VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
        `, userID, productID, quantity)
		if err != nil {
			return fmt.Errorf("failed to insert new product into cart: %w", err)
		}
	case err == nil:
		_, err = tx.ExecContext(ctx, `
            UPDATE carts_products
            SET quantity = ?, modified_at = CURRENT_TIMESTAMP
            WHERE user_id = ? AND product_id = ?;
        `, currentQty+quantity, userID, productID)
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
	default:
		return fmt.Errorf("failed to query cart: %w", err)
	}

	return tx.Commit()
}

func (repo *MySQLCartRepository) GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error) {
	rows, err := repo.DB.QueryContext(ctx, `
        SELECT cp.id, cp.user_id, cp.product_id, cp.quantity, cp.created_at, cp.modified_at,
               p.id, p.name, p.description, p.price, p.stock
        FROM carts_products cp
        JOIN products p ON cp.product_id = p.id
        WHERE cp.user_id = ?;
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	defer rows.Close()

	var items []*entities.CartItem
	for rows.Next() {
		var item entities.CartItem
		var product entities.Product
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.ModifiedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock,
		); err != nil {
			return nil, err
		}
		item.Product = &product
		items = append(items, &item)
	}

	return items, rows.Err()
}

func (repo *MySQLCartRepository) UpdateCartItem(ctx context.Context, userID, productID int64, quantity int) error {
	_, err := repo.DB.ExecContext(ctx, `
        UPDATE carts_products
        SET quantity = ?, modified_at = CURRENT_TIMESTAMP
        WHERE user_id = ? AND product_id = ?;
    `, quantity, userID, productID)
	return err
}

func (repo *MySQLCartRepository) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	_, err := repo.DB.ExecContext(ctx, `
        DELETE FROM carts_products
        WHERE user_id = ? AND product_id = ?;
    `, userID, productID)
	return err
}

func (repo *MySQLCartRepository) ClearCart(ctx context.Context, userID int64) error {
	_, err := repo.DB.ExecContext(ctx, `DELETE FROM carts_products WHERE user_id = ?;`, userID)
	return err
}
