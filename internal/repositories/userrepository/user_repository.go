package userrepository

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"database/sql"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlErrConflict uint16 = 1062
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetAll users from the database, or an error if any occurs
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User

	rows, err := r.DB.Query("SELECT id, email, password, phone, is_adm FROM USERS")
	if err != nil {
		return nil, core.ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
			return nil, core.ErrInternal
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, core.ErrInternal
	}

	if users == nil {
		users = []models.User{}
	}

	return users, nil
}

// Add a user to the database, or an error if any occurs
func (r *UserRepository) Add(user models.RegisterRequest) (*models.UserResponse, error) {
	const query = "INSERT INTO users (email, password, phone, is_adm) VALUES (?, ?, ?, ?)"

	res, err := r.DB.Exec(query, user.Email, user.Password, user.Phone, user.IsAdm)
	if err != nil {
		log.Printf("err: %v", err)

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlErrConflict {
			return nil, core.ErrConflict
		}

		return nil, core.ErrInternal
	}

	id, _ := res.LastInsertId()
	return &models.UserResponse{ID: id}, nil
}

// Delete a user from the database with the given userId. Return an error if any occurs
func (r *UserRepository) Delete(userId int64) error {
	const query = "DELETE FROM users WHERE id = ?"

	_, err := r.DB.Exec(query, userId)

	if err != nil {
		return err
	}

	return nil
}

// FindById returns the user with the given userId, or an error if any occurs
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE email = ?"

	row := r.DB.QueryRow(query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternal
	}

	return &user, nil
}

// FindById returns the user with the given userId in the database, or an error if any occurs
func (r *UserRepository) FindById(userId int64) (*models.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE id = ?"

	row := r.DB.QueryRow(query, userId)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternal
	}

	return &user, nil
}

func (r *UserRepository) Update(user models.UserRequest, userId int64) (*models.UserResponse, error) {
	existing, err := r.FindById(userId)

	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, core.ErrNotFound
	}

	if user.Email != "" {
		existing.Email = user.Email
	}

	if user.Phone != "" {
		existing.Phone = user.Phone
	}

	existing.IsAdm = user.IsAdm

	const query = "UPDATE users SET email = ?, password = ?, phone = ?, id_adm = ? WHERE id = ?"
	_, err = r.DB.Exec(query, existing.Email, existing.Password, existing.Phone, existing.IsAdm, userId)

	log.Printf("Err: %v", err)

	return &models.UserResponse{ID: userId}, nil
}
