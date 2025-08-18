package user

import (
	"cachacariaapi/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAll() []models.User {
	var users []models.User

	req, err := r.DB.Query("SELECT id, name, email, password, phone, is_adm FROM USERS")

	if err != nil {
		log.Fatal(err)
	}

	defer req.Close()

	log.Printf("Request: %v", req)

	for req.Next() {
		var user models.User

		if err := req.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}

func (r *UserRepository) Add(user models.AddUserRequest) (*models.AddUserResponse, error) {
	const query = "INSERT INTO USERS (name, email, password, phone, is_adm) VALUES (?, ?, ?, ?, ?)"

	//id := uuid.New()
	res, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.IsAdm)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Result: %v", res)

	id, _ := res.LastInsertId()
	return &models.AddUserResponse{ID: id}, nil
}
