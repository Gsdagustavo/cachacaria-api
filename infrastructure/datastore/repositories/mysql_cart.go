package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/infrastructure/datastore"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
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

	res, err := tx.ExecContext(ctx, `
		INSERT INTO orders (user_id)
		VALUES (?);
	`, userID)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to create order"), err)
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		return errors.Join(fmt.Errorf("failed to retrieve order ID"), err)
	}

	for _, it := range items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES (?, ?, ?);
		`, orderID, it.ProductID, it.Quantity)
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

func (repo *MySQLCartRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*entities.Order, error) {
	rows, err := repo.DB.QueryContext(ctx, `
        SELECT
            o.id AS order_id,
            o.user_id,
            o.created_at AS order_created_at,
            o.modified_at AS order_modified_at,

            oi.id AS order_item_id,
            oi.product_id,
            oi.quantity,
            oi.created_at AS item_created_at,
            oi.modified_at AS item_modified_at,

            p.id        AS product_id,
            p.name      AS product_name,
            p.description AS product_description,
            p.price     AS product_price,
            p.stock     AS product_stock,
            (
                SELECT GROUP_CONCAT(pp.filename)
                FROM products_photos pp
                WHERE pp.product_id = p.id
            ) AS photos
        FROM orders o
        JOIN order_items oi ON oi.order_id = o.id
        JOIN products p ON p.id = oi.product_id
        WHERE o.user_id = ?
        ORDER BY o.created_at DESC, oi.id ASC;
    `, userID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to query orders"), err)
	}
	defer rows.Close()

	ordersMap := make(map[int64]*entities.Order)

	for rows.Next() {
		var (
			orderID       int64
			userIDDB      int64
			orderCreated  time.Time
			orderModified time.Time

			itemID       int64
			productID    int64
			quantity     int
			itemCreated  time.Time
			itemModified time.Time

			pID          int64
			productName  sql.NullString
			productDesc  sql.NullString
			productPrice float64
			productStock int
			photosStr    sql.NullString
		)

		if err = rows.Scan(
			&orderID,
			&userIDDB,
			&orderCreated,
			&orderModified,

			&itemID,
			&productID,
			&quantity,
			&itemCreated,
			&itemModified,

			&pID,
			&productName,
			&productDesc,
			&productPrice,
			&productStock,
			&photosStr,
		); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to scan orders rows"), err)
		}

		// cria o pedido se ainda não existir
		if _, ok := ordersMap[orderID]; !ok {
			ordersMap[orderID] = &entities.Order{
				ID:          orderID,
				UserID:      userIDDB,
				Status:      "completed", // ajuste se tiver campo real
				TotalAmount: 0,
				CreatedAt:   orderCreated,
				ModifiedAt:  orderModified,
				Items:       []entities.OrderItem{},
			}
		}

		// monta fotos e produto
		photos := []string{}
		if photosStr.Valid && photosStr.String != "" {
			// GROUP_CONCAT usa ',' por padrão
			photos = strings.Split(photosStr.String, ",")
		}

		prod := &entities.Product{
			ID:          pID,
			Name:        "",
			Description: "",
			Photos:      photos,
			Price:       float32(productPrice),
			Stock:       productStock,
		}
		if productName.Valid {
			prod.Name = productName.String
		}
		if productDesc.Valid {
			prod.Description = productDesc.String
		}

		item := entities.OrderItem{
			ID:         itemID,
			OrderID:    orderID,
			ProductID:  productID,
			Product:    prod,
			Quantity:   quantity,
			Price:      productPrice,
			CreatedAt:  itemCreated,
			ModifiedAt: itemModified,
		}

		o := ordersMap[orderID]
		o.Items = append(o.Items, item)
		o.TotalAmount += float64(item.Quantity) * productPrice
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(fmt.Errorf("rows iteration error"), err)
	}

	// transforma map em slice ordenado (já estava ordenado pela query, mas map perde ordem)
	result := make([]*entities.Order, 0, len(ordersMap))
	for _, o := range ordersMap {
		result = append(result, o)
	}

	// opcional: ordenar por CreatedAt desc (caso queira garantir ordem)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}
