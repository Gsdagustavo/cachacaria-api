package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/infrastructure/datastore"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type MySQLCartRepository struct {
	DB *sql.DB
}

func NewMySQLCartRepository(db *sql.DB) repositories.CartRepository {
	return &MySQLCartRepository{DB: db}
}

func (repo *MySQLCartRepository) AddToCart(ctx context.Context, userID, productID int64, quantity int) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to begin transaction"), err)
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
			return errors.Join(fmt.Errorf("failed to insert new product into cart"), err)
		}
	case err == nil:
		_, err = tx.ExecContext(ctx, `
            UPDATE carts_products
            SET quantity = ?, modified_at = CURRENT_TIMESTAMP
            WHERE user_id = ? AND product_id = ?;
        `, currentQty+quantity, userID, productID)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to update product"), err)
		}
	default:
		return errors.Join(fmt.Errorf("failed to query for products"), err)
	}

	return tx.Commit()
}

func (repo *MySQLCartRepository) GetCartItems(ctx context.Context, userID int64) ([]*entities.CartItem, error) {
	rows, err := repo.DB.QueryContext(ctx, `
        SELECT 
            cp.id, cp.user_id, cp.product_id, cp.quantity, cp.created_at, cp.modified_at,
            p.id, p.name, p.description, p.price, p.stock,
            GROUP_CONCAT(pp.filename) AS photos
        FROM carts_products cp
        JOIN products p ON cp.product_id = p.id
        LEFT JOIN products_photos pp ON pp.product_id = p.id
        WHERE cp.user_id = ?
        GROUP BY cp.id, p.id;
    `, userID)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to query cart items"), err)
	}
	defer rows.Close()

	var items []*entities.CartItem

	for rows.Next() {
		var item entities.CartItem
		var product entities.Product
		var photosStr sql.NullString

		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.ModifiedAt, &product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &photosStr)

		if err != nil {
			return nil, errors.Join(fmt.Errorf("failed to scan cart"), err)
		}

		if photosStr.Valid {
			product.Photos = strings.Split(photosStr.String, ",")
		} else {
			product.Photos = []string{}
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
	if err != nil {
		return errors.Join(fmt.Errorf("failed to execute cart item update"), err)
	}

	return nil
}

func (repo *MySQLCartRepository) DeleteCartItem(ctx context.Context, userID, productID int64) error {
	_, err := repo.DB.ExecContext(ctx, `
        DELETE FROM carts_products
        WHERE user_id = ? AND product_id = ?;
    `, userID, productID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to execute cart item deletion"), err)
	}

	return nil
}

func (repo *MySQLCartRepository) ClearCart(ctx context.Context, userID int64) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to begin transaction"), err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
        SELECT product_id, quantity
        FROM carts_products
        WHERE user_id = ?;
    `, userID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to query cart before checkout"), err)
	}
	defer rows.Close()

	items := []struct {
		ProductID int64
		Quantity  int
	}{}

	for rows.Next() {
		var item struct {
			ProductID int64
			Quantity  int
		}
		if err = rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return errors.Join(fmt.Errorf("failed to scan cart item"), err)
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	for _, it := range items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO orders (user_id, product_id, quantity, created_at)
            VALUES (?, ?, ?, CURRENT_TIMESTAMP);
        `, userID, it.ProductID, it.Quantity)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to insert order item"), err)
		}
	}

	_, err = tx.ExecContext(ctx, `
        DELETE FROM carts_products
        WHERE user_id = ?;
    `, userID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to clear cart"), err)
	}

	return tx.Commit()
}
