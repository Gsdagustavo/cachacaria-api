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
	return &MySQLCartRepository{
		DB: db,
	}
}

func (repo *MySQLCartRepository) GetCartByUserID(userID int64) ([]entities.CartItem, error) {
	query := `
        SELECT 
            cp.id,
            cp.user_id,
            cp.product_id,
            cp.quantity,
            cp.created_at,
            cp.modified_at,
            p.name,
            p.description,
            p.price
        FROM carts_products cp
        JOIN products p ON cp.product_id = p.id
        WHERE cp.user_id = ?;
    `

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to fetch cart for user"), err)
	}
	defer rows.Close()

	var items []entities.CartItem
	for rows.Next() {
		var item entities.CartItem
		var product entities.Product
		err = rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.ModifiedAt, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}
		item.Product = &product
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *MySQLCartRepository) AddProductToCart(userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity to add must be greater than zero")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return errors.Join(fmt.Errorf("failed to begin transaction"), err)
	}
	defer func() {
		_ = tx.Rollback() // rollback if not committed
	}()

	// Try to get the current quantity
	var currentQty int
	querySelect := `
        SELECT quantity
        FROM carts_products
        WHERE user_id = ? AND product_id = ?
        FOR UPDATE;
    `
	err = tx.QueryRow(querySelect, userID, productID).Scan(&currentQty)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Product not in cart → insert new entry
			queryInsert := `
                INSERT INTO carts_products (user_id, product_id, quantity)
                VALUES (?, ?, ?);
            `
			_, err = tx.Exec(queryInsert, userID, productID, quantity)
			if err != nil {
				return errors.Join(fmt.Errorf("failed to insert new product into cart"), err)
			}
		} else {
			return errors.Join(fmt.Errorf("failed to check existing product quantity"), err)
		}
	} else {
		// Product exists → update its quantity
		newQty := currentQty + quantity
		queryUpdate := `
            UPDATE carts_products
            SET quantity = ?, modified_at = CURRENT_TIMESTAMP
            WHERE user_id = ? AND product_id = ?;
        `
		_, err = tx.Exec(queryUpdate, newQty, userID, productID)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to update product quantity"), err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(fmt.Errorf("failed to commit transaction"), err)
	}

	return nil
}

func (repo *MySQLCartRepository) DecreaseProductQuantity(userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity to decrease must be greater than zero")
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return errors.Join(fmt.Errorf("failed to begin transaction"), err)
	}
	defer func() {
		_ = tx.Rollback() // safe rollback if not committed
	}()

	// Get current quantity
	var currentQty int
	querySelect := `
        SELECT quantity
        FROM carts_products
        WHERE user_id = ? AND product_id = ?
        FOR UPDATE;
    `
	err = tx.QueryRow(querySelect, userID, productID).Scan(&currentQty)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("product not found in cart")
		}
		return errors.Join(fmt.Errorf("failed to fetch current quantity"), err)
	}

	// Calculate new quantity
	newQty := currentQty - quantity
	if newQty <= 0 {
		// Remove product entirely if quantity goes to zero or below
		queryDelete := `
            DELETE FROM carts_products
            WHERE user_id = ? AND product_id = ?;
        `
		_, err = tx.Exec(queryDelete, userID, productID)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to delete product from cart"), err)
		}
	} else {
		// Just update the quantity
		queryUpdate := `
            UPDATE carts_products
            SET quantity = ?, modified_at = CURRENT_TIMESTAMP
            WHERE user_id = ? AND product_id = ?;
        `
		_, err = tx.Exec(queryUpdate, newQty, userID, productID)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to update product quantity"), err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(fmt.Errorf("failed to commit transaction"), err)
	}

	return nil
}

func (repo *MySQLCartRepository) RemoveProductFromCart(userID, productID int64) error {
	query := `DELETE FROM carts_products WHERE user_id = ? AND product_id = ?`
	_, err := repo.DB.Exec(query, userID, productID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to remove product from cart"), err)
	}

	return nil
}
