package persistence

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/interfaces/http/core"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MySQLProductRepository struct {
	DB *sql.DB
}

func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
	return &MySQLProductRepository{DB: db}
}

func (r *MySQLProductRepository) Add(product entities.AddProductRequest, photos []*multipart.FileHeader) (*entities.AddProductResponse, error) {
	query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"

	res, err := r.DB.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		log.Printf("MySQLProductsRepository.Add Error: %v", err)
		return nil, err
	}

	id, _ := res.LastInsertId()

	for _, fileheader := range photos {
		src, err := fileheader.Open()
		if err != nil {
			log.Printf("MySQLProductsRepository.Add Error: %v", err)
			return nil, err
		}
		defer src.Close()

		filename := fmt.Sprintf("\"product_%d_%d%s", id, time.Now().UnixNano(), filepath.Ext(fileheader.Filename))
		filePath := filepath.Join("/app/images", filename)

		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("MySQLProductsRepository.Add Error: %v", err)
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("MySQLProductsRepository.Add Error: %v", err)
			return nil, err
		}

		const query = "INSERT INTO products_photos (product_id, filename) VALUES (?, ?)"
		_, err = r.DB.Exec(query, id, filename)

		if err != nil {
			log.Printf("MySQLProductsRepository.Add Error: %v", err)
			return nil, err
		}
	}

	return &entities.AddProductResponse{ID: id}, nil
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

	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, core.ErrInternal
		}
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

func decodeBase64Image(data string) ([]byte, error) {
	// Check for the "data:image/...;base64," prefix
	if idx := strings.Index(data, ","); idx != -1 {
		data = data[idx+1:] // remove the prefix
	}

	// Remove any whitespace/newlines
	data = strings.TrimSpace(data)

	// Decode the Base64 string
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.New("error decoding base64 image: " + err.Error())
	}

	return imgBytes, nil
}
