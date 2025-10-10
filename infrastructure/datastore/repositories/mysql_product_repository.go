package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/interfaces/http/core"
	"database/sql"
	"errors"
	"log/slog"
)

type MySQLProductRepository struct {
	DB *sql.DB
}

func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
	return &MySQLProductRepository{DB: db}
}

func (r *MySQLProductRepository) AddProduct(product entities.AddProductRequest) (*entities.AddProductResponse, error) {
	query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"

	res, err := r.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		slog.Error("[MySQLProductsRepository.AddProduct] error adding product", "error", err.Error())
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &entities.AddProductResponse{ID: id}, nil
}

func (r *MySQLProductRepository) AddProductPhotos(productID int64, filenames []string) error {
	for _, f := range filenames {
		_, err := r.DB.Exec("INSERT INTO products_photos (product_id, filename) VALUES (?, ?)", productID, f)
		if err != nil {
			slog.Error("[MySQLProductRepository.AddProductPhotos] error adding product photos", "error", err.Error())
			return err
		}
	}
	return nil
}

func (r *MySQLProductRepository) GetAll(limit, offset int) ([]entities.Product, error) {
	var products []entities.Product

	const query = "SELECT id, name, description, price, stock FROM products ORDER BY id LIMIT ? OFFSET ?"
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		slog.Error("[MySQLProductRepository.getAll] error getting all products", "error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return products, nil
		}
		return nil, core.ErrInternal
	}

	defer rows.Close()

	/// Products
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			slog.Error("[MySQLProductRepository.getAll] error scanning products rows", "error", err.Error())
			return nil, core.ErrInternal
		}

		var photos []string
		const query = "SELECT id, filename FROM products_photos WHERE product_id = ?"
		photoRows, err := r.DB.Query(query, product.ID)
		if err != nil {
			slog.Warn("[MySQLProductRepository.getAll] warning on get product select", "warning", err.Error(), "query", query)
			continue
		}

		for photoRows.Next() {
			var photoID int64
			var filename string
			if err := photoRows.Scan(&photoID, &filename); err != nil {
				slog.Error("[MySQLProductRepository.getAll] error scanning photo rows", "error", err.Error())
				return nil, core.ErrInternal
			}

			photos = append(photos, filename)
		}

		product.Photos = photos
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		slog.Error("[MySQLProductRepository.getAll] error on get product photos", "error", err.Error())
		return nil, core.ErrInternal
	}

	if products == nil {
		products = []entities.Product{}
	}

	return products, nil
}

func (r *MySQLProductRepository) GetProduct(id int64) (*entities.Product, error) {
	const query = "SELECT id, name, description, price, stock FROM products WHERE id = ?"
	row := r.DB.QueryRow(query, id)

	var product entities.Product

	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
		slog.Error("[MySQLProductRepository.get] error on scanning proudcts", "error", err.Error(), "query", query)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternal
	}

	const photoQuery = "SELECT id, filename FROM products_photos WHERE product_id = ?"
	photoRows, err := r.DB.Query(photoQuery, product.ID)
	if err != nil {
		slog.Error("[MySQLProductRepository.get] error on deleting product", "error", err.Error(), "query", query)
	} else {
		defer photoRows.Close()
		photos := make([]string, 0)
		for photoRows.Next() {
			var photoID int64
			var filename string
			if err := photoRows.Scan(&photoID, &filename); err != nil {
				slog.Error("[MySQLProductRepository.get] error on scanning product photos row", "error", err.Error(), "query", query)
				return nil, core.ErrInternal
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

func (r *MySQLProductRepository) DeleteProduct(id int64) (*entities.DeleteProductResponse, error) {
	const query = "DELETE FROM products WHERE id = ?"

	txn, err := r.DB.Begin()
	if err != nil {
		slog.Error("[MySQLProductRepository.delete] error on deleting product", "error", err.Error(), "query", query)
		return nil, err
	}

	res, err := txn.Exec(query, id)
	if err != nil {
		slog.Error("[MySQLProductRepository.delete] error on deleting product", "error", err.Error(), "query", query)
		txn.Rollback()
		return nil, err
	}

	if rows, err := res.RowsAffected(); err != nil || rows != 1 {
		slog.Error("[MySQLProductRepository.delete] error on deleting product", "error", err.Error(), "query", query)
		txn.Rollback()
		return nil, err
	}

	if err := txn.Commit(); err != nil {
		slog.Error("[MySQLProductRepository.delete] error on deleting product", "error", err.Error(), "query", query)
		txn.Rollback()
		return nil, err
	}

	return &entities.DeleteProductResponse{ID: id}, nil
}

func (r *MySQLProductRepository) UpdateProduct(id int64, product entities.UpdateProductRequest) (*entities.UpdateProductResponse, error) {
	const query = "UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?"
	res, err := r.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock, id)
	if err != nil {
		slog.Error("[MySQLProductRepository.update] error on update product", "error", err.Error(), "query", query)
		return nil, err
	}

	if rows, err := res.RowsAffected(); err != nil || rows != 1 {
		slog.Error("[MySQLProductRepository.update] no rows affected on update product", "error", err.Error(), "query", query)
		return nil, err
	}

	return &entities.UpdateProductResponse{ID: id}, nil
}
