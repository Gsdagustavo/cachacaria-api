package persistence

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/interfaces/http/core"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

type MySQLUserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{DB: db}
}

// GetAll users from the database, or an error if any occurs
func (r *MySQLUserRepository) GetAll() ([]entities.User, error) {
	var users []entities.User

	rows, err := r.DB.Query("SELECT id, email, password, phone, is_adm FROM users")
	if err != nil {
		slog.Error("[MySQLUserRepository.GetAll] error getting users", "error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}

		return nil, core.ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
			slog.Error("[MySQLUserRepository.GetAll] error scanning users row", "error", err.Error())
			return nil, core.ErrInternal
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("[MySQLUserRepository.GetAll] error getting users", "error", err.Error())
		return nil, core.ErrInternal
	}

	if users == nil {
		users = []entities.User{}
	}

	return users, nil
}

// Add a user to the database. Returns a UserResponse or an error if any occurs
func (r *MySQLUserRepository) Add(user entities.RegisterRequest) (*entities.UserResponse, error) {
	const query = "INSERT INTO users (email, password, phone, is_adm) VALUES (?, ?, ?, ?)"

	res, err := r.DB.Exec(query, user.Email, user.Password, user.Phone, user.IsAdm)
	if err != nil {
		slog.Error("[MySQLUserRepository.Add] error adding user", "error", err.Error(), "query", query)

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlErrConflict {
			slog.Warn("[MySQLUserRepository.Add] warning adding user", "error", err.Error(), "cause", "mysql conflict", "query", query)
			return nil, core.ErrConflict
		}

		return nil, core.ErrInternal
	}

	id, _ := res.LastInsertId()

	slog.Info(fmt.Sprintf("[MySQLUserRepository.Add] user with id %v added successfully", id))

	return &entities.UserResponse{ID: id}, nil
}

// Delete a user from the database with the given userId. Return an error if any occurs
func (r *MySQLUserRepository) Delete(userId int64) error {
	const query = "DELETE FROM users WHERE id = ?"

	_, err := r.DB.Exec(query, userId)

	if err != nil {
		slog.Error("[MySQLUserRepository.Delete] error deleting user", "error", err.Error(), "query", query)
		return err
	}

	slog.Info(fmt.Sprintf("[MySQLUserRepository.Delete] user with id %v deleted successfully", userId))

	return nil
}

// FindByEmail returns the user with the given email, or an error if any occurs
func (r *MySQLUserRepository) FindByEmail(email string) (*entities.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE email = ?"

	row := r.DB.QueryRow(query, email)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		slog.Error("[MySQLUserRepository.FindByEmail] error scanning user rows", "error", err.Error(), "query", query)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternal
	}

	return &user, nil
}

// FindById returns the user with the given userId in the database, or an error if any occurs
func (r *MySQLUserRepository) FindById(userId int64) (*entities.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE id = ?"

	row := r.DB.QueryRow(query, userId)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		slog.Error("[MySQLUserRepository.FindById] error scanning user rows", "error", err.Error(), "query", query)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		return nil, core.ErrInternal
	}

	return &user, nil
}

func (r *MySQLUserRepository) Update(user entities.UserRequest, userId int64) (*entities.UserResponse, error) {
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
	if err != nil {
		slog.Error("[MySQLUserRepository.Update] error updating user", "error", err.Error(), "query", query)
		return nil, err
	}

	slog.Info(fmt.Sprintf("[MySQLUserRepository.Update] user with id %v updated successfully", existing.ID))

	return &entities.UserResponse{ID: userId}, nil
}
