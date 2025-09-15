package persistence

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/interfaces/http/core"
	"database/sql"
	"errors"
	"log"
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
		log.Printf("MySQLProductsRepository.AddProduct Error: %v", err)
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &entities.AddProductResponse{ID: id}, nil
}

func (r *MySQLProductRepository) AddProductPhotos(productID int64, filenames []string) error {
	for _, f := range filenames {
		_, err := r.DB.Exec("INSERT INTO products_photos (product_id, filename) VALUES (?, ?)", productID, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLProductRepository) GetAll() ([]entities.Product, error) {
	var products []entities.Product

	rows, err := r.DB.Query("SELECT id, name, description, price, stock FROM products")
	if err != nil {
		log.Printf("MySQLUserRepository.GetAll Error: %v", err)

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
			return nil, core.ErrInternal
		}

		var photos []string
		const query = "SELECT id, filename FROM products_photos WHERE product_id = ?"
		photoRows, err := r.DB.Query(query, product.ID)
		if err != nil {
			log.Printf("MySQLUserRepository.GetAll Error: %v", err)
			continue
		}

		for photoRows.Next() {
			var photoID int64
			var filename string
			if err := photoRows.Scan(&photoID, &filename); err != nil {
				return nil, core.ErrInternal
			}

			log.Printf("filename: %v", filename)

			photos = append(photos, filename)
		}

		log.Printf("Products photos: %v", photos)

		product.Photos = photos
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, core.ErrInternal
	}

	if products == nil {
		products = []entities.Product{}
	}

	return products, nil
}
