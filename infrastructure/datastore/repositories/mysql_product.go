package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/infrastructure/datastore"
	"database/sql"
	"errors"
	"fmt"
)

type MySQLProductRepository struct {
	DB *sql.DB
}

func NewMySQLProductRepository(db *sql.DB) repositories.ProductRepository {
	return &MySQLProductRepository{DB: db}
}

func (r *MySQLProductRepository) AddProduct(product entities.AddProductRequest) (int64, error) {
	query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"

	res, err := r.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return -1, errors.Join(fmt.Errorf("failed to insert product"), err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (r *MySQLProductRepository) AddProductPhotos(productID int64, filenames []string) error {
	for _, f := range filenames {
		_, err := r.DB.Exec("INSERT INTO products_photos (product_id, filename) VALUES (?, ?)", productID, f)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to insert product"), err)
		}
	}
	return nil
}

func (r *MySQLProductRepository) GetAll(limit, offset int) ([]entities.Product, error) {
	products := make([]entities.Product, 0)

	/// Products
	const query = "SELECT id, name, description, price, stock FROM products ORDER BY id LIMIT ? OFFSET ? WHERE status_code != 1"
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return products, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query products"), err)
	}

	defer rows.Close()

	for rows.Next() {
		var product entities.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to scan product"), err)
		}

		// Photos
		var photos []string
		const query = "SELECT id, filename FROM products_photos WHERE product_id = ?"
		photoRows, err := r.DB.Query(query, product.ID)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("failed to query products images"), err)
		}

		for photoRows.Next() {
			var photoID int64
			var filename string
			if err = photoRows.Scan(&photoID, &filename); err != nil {
				return nil, errors.Join(fmt.Errorf("failed to scan product"), err)
			}

			photos = append(photos, filename)
		}

		product.Photos = photos
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to scan products"), err)
	}

	return products, nil
}

func (r *MySQLProductRepository) GetProduct(id int64) (*entities.Product, error) {
	const query = "SELECT id, name, description, price, stock FROM products WHERE id = ? and status_code != 1"
	row := r.DB.QueryRow(query, id)

	var product entities.Product

	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(errors.New("failed to scan product"), err)
	}

	const photoQuery = "SELECT id, filename FROM products_photos WHERE product_id = ?"
	photoRows, err := r.DB.Query(photoQuery, product.ID)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to query product photos"), err)
	} else {
		defer photoRows.Close()
		photos := make([]string, 0)
		for photoRows.Next() {
			var photoID int64
			var filename string
			if err = photoRows.Scan(&photoID, &filename); err != nil {
				return nil, errors.Join(fmt.Errorf("failed to scan product photo"), err)
			}

			photos = append(photos, filename)
		}

		product.Photos = photos
	}

	if product.Photos == nil {
		product.Photos = make([]string, 0)
	}

	return &product, nil
}

func (r *MySQLProductRepository) DeleteProduct(id int64) error {
	const query = `UPDATE products SET status_code = 1 WHERE id = ?`

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to update product"), err)
	}

	return nil
}

func (r *MySQLProductRepository) UpdateProduct(id int64, product entities.UpdateProductRequest) error {
	const query = "UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?"
	_, err := r.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock, id)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to update product"), err)
	}

	return nil
}

func (r *MySQLProductRepository) DecrementStock(id int64, quantity int) error {
	const query = "UPDATE products SET stock = stock - ? WHERE id = ?"

	_, err := r.DB.Exec(query, quantity, id)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to update product"), err)
	}

	return nil
}
