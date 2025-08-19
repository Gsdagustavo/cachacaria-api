package userrepository

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

	res, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.IsAdm)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Result: %v", res)

	id, _ := res.LastInsertId()
	return &models.AddUserResponse{ID: id}, nil
}

func (r *UserRepository) Delete(userId int64) error {
	_, err := r.DB.Exec("DELETE FROM USERS WHERE ID = ?", userId)

	if err != nil {
		log.Printf("Error on deleting user: %v", err)
		return err
	}

	log.Printf("User with ID %v removed successfully", userId)

	return nil
}

func (r *UserRepository) FindById(userId int64) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, name, email, password, phone, is_adm FROM USERS WHERE ID = ?", userId)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %v not found", userId)
		}
		return nil, err
	}

	return &user, nil
}
