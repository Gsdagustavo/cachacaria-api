package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/util"
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

func (repo *MySQLCartRepository) GetCartItems(ctx context.Context, userID int64, baseURL string) ([]*entities.CartItem, error) {
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

		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity,
			&item.CreatedAt, &item.ModifiedAt,
			&product.ID, &product.Name, &product.Description,
			&product.Price, &product.Stock,
			&photosStr,
		)

		if err != nil {
			return nil, errors.Join(fmt.Errorf("failed to scan cart"), err)
		}

		// Converter GROUP_CONCAT -> []string
		if photosStr.Valid {
			filenames := strings.Split(photosStr.String, ",")

			product.Photos = make([]string, len(filenames))
			for i, filename := range filenames {
				product.Photos[i] = util.GetProductImageURL(filename, baseURL)
			}
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
	_, err := repo.DB.ExecContext(ctx, `DELETE FROM carts_products WHERE user_id = ?;`, userID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to execute clear cart query"), err)
	}

	return nil
}
