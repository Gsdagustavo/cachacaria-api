package repositories

import (
	"cachacariaapi/domain/entities"
	"database/sql"
	"errors"
	"fmt"
)

type MySQLCartRepository struct {
	DB *sql.DB
}

func NewMySQLCartRepository(db *sql.DB) *MySQLCartRepository {
	return &MySQLCartRepository{DB: db}
}

// Get all items in the user's cart, with product data joined in.
func (repo *MySQLCartRepository) GetCartByUserID(userID int64) ([]entities.CartItem, error) {
	query := `
		SELECT 
			cp.id,
			cp.user_id,
			cp.product_id,
			cp.quantity,
			cp.created_at,
			cp.modified_at,
			p.id,
			p.name,
			p.description,
			p.price,
			p.stock
		FROM carts_products cp
		JOIN products p ON cp.product_id = p.id
		WHERE cp.user_id = ?;
	`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart for user: %w", err)
	}
	defer rows.Close()

	var items []entities.CartItem
	for rows.Next() {
		var item entities.CartItem
		var product entities.Product

		err = rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.ModifiedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock,
		)
		if err != nil {
			return nil, err
		}

		item.Product = &product
		items = append(items, item)
	}

	return items, rows.Err()
}

// Add or increase a product’s quantity in the cart.
func (repo *MySQLCartRepository) AddProductToCart(userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRow(`
		SELECT quantity
		FROM carts_products
		WHERE user_id = ? AND product_id = ? FOR UPDATE;
	`, userID, productID).Scan(&currentQty)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = tx.Exec(`
			INSERT INTO carts_products (user_id, product_id, quantity, created_at, modified_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
		`, userID, productID, quantity)
		if err != nil {
			return fmt.Errorf("failed to insert new product into cart: %w", err)
		}
	case err == nil:
		_, err = tx.Exec(`
			UPDATE carts_products
			SET quantity = ?, modified_at = CURRENT_TIMESTAMP
			WHERE user_id = ? AND product_id = ?;
		`, currentQty+quantity, userID, productID)
		if err != nil {
			return fmt.Errorf("failed to update existing product quantity: %w", err)
		}
	default:
		return fmt.Errorf("failed to query cart: %w", err)
	}

	return tx.Commit()
}

// Decrease quantity or remove product entirely.
func (repo *MySQLCartRepository) DecreaseProductQuantity(userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRow(`
		SELECT quantity
		FROM carts_products
		WHERE user_id = ? AND product_id = ? FOR UPDATE;
	`, userID, productID).Scan(&currentQty)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("product not found in cart")
	}
	if err != nil {
		return fmt.Errorf("failed to query cart: %w", err)
	}

	if currentQty <= quantity {
		_, err = tx.Exec(`
			DELETE FROM carts_products
			WHERE user_id = ? AND product_id = ?;
		`, userID, productID)
	} else {
		_, err = tx.Exec(`
			UPDATE carts_products
			SET quantity = ?, modified_at = CURRENT_TIMESTAMP
			WHERE user_id = ? AND product_id = ?;
		`, currentQty-quantity, userID, productID)
	}
	if err != nil {
		return fmt.Errorf("failed to update or delete product: %w", err)
	}

	return tx.Commit()
}

// Remove product entirely.
func (repo *MySQLCartRepository) RemoveProductFromCart(userID, productID int64) error {
	_, err := repo.DB.Exec(`
		DELETE FROM carts_products
		WHERE user_id = ? AND product_id = ?;
	`, userID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove product from cart: %w", err)
	}
	return nil
}

// Clear all cart items (useful for order checkout).
func (repo *MySQLCartRepository) ClearCart(userID int64) error {
	_, err := repo.DB.Exec(`DELETE FROM carts_products WHERE user_id = ?;`, userID)
	return err
}
