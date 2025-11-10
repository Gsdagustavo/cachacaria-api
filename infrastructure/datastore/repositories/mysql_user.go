package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/infrastructure/datastore"
	"database/sql"
	"errors"
	"fmt"
)

type MySQLUserRepository struct {
	DB *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) repositories.UserRepository {
	return &MySQLUserRepository{DB: db}
}

// GetAll users from the database, or an error if any occurs
func (r *MySQLUserRepository) GetAll() ([]entities.User, error) {
	var users []entities.User

	rows, err := r.DB.Query("SELECT id, email, password, phone, is_adm FROM users")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query all users"), err)
	}

	defer rows.Close()

	for rows.Next() {
		var user entities.User
		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to scan users row"), err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to scan users row"), err)
	}

	if users == nil {
		users = []entities.User{}
	}

	return users, nil
}

// Add a user to the database. Returns a UserResponse or an error if any occurs
func (r *MySQLUserRepository) Add(user entities.User) error {
	const query = "INSERT INTO users (uuid, email, password, phone, is_adm) VALUES (?, ?, ?, ?)"

	_, err := r.DB.Exec(query, user.Email, user.Password, user.Phone, user.IsAdm)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to insert user"), err)
	}

	return nil
}

// Delete a user from the database with the given userId. Return an error if any occurs
func (r *MySQLUserRepository) Delete(userId int64) error {
	const query = "DELETE FROM users WHERE id = ?"

	_, err := r.DB.Exec(query, userId)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to delete user"), err)
	}

	return nil
}

// FindByEmail returns the user with the given email, or an error if any occurs
func (r *MySQLUserRepository) FindByEmail(email string) (*entities.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE email = ?"

	row := r.DB.QueryRow(query, email)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query user"), err)
	}

	return &user, nil
}

// FindById returns the user with the given userId in the database, or an error if any occur
func (r *MySQLUserRepository) FindById(userId int64) (*entities.User, error) {
	const query = "SELECT id, email, password, phone, is_adm FROM users WHERE id = ?"

	row := r.DB.QueryRow(query, userId)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query user"), err)
	}

	return &user, nil
}

func (r *MySQLUserRepository) Update(user entities.User, userId int64) error {
	existing, err := r.FindById(userId)

	if err != nil {
		return nil
	}

	if existing == nil {
		return nil
	}

	if user.Email != "" {
		existing.Email = user.Email
	}

	if user.Phone != "" {
		existing.Phone = user.Phone
	}

	existing.IsAdm = user.IsAdm

	const query = "UPDATE users SET email = ?, password = ?, phone = ?, id_adm = ? WHERE id = ?"
	_, err = r.DB.Exec(
		query,
		existing.Email,
		existing.Password,
		existing.Phone,
		existing.IsAdm,
		userId,
	)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to update user"), err)
	}

	return nil
}
