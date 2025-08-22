package userrepository

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetAll users from the database, or an error if any occurs
func (r *UserRepository) GetAll() ([]models.User, error) {
	// instantiate a new slice of users
	var users []models.User

	// query through users table
	rows, err := r.DB.Query("SELECT id, name, email, password, phone, is_adm FROM USERS")

	if err != nil {
		return nil, core.ErrInternal
	}
	defer rows.Close()

	// read users
	for rows.Next() {
		// instantiate a new user
		var user models.User

		// scan the user from the given row
		// if any error occurs during the scanning process, return err internal
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {
			return nil, core.ErrInternal
		}

		// append the new user in the users list
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, core.ErrNotFound
	}

	return users, nil
}

// Add a user to the database, or an error if any occurs
func (r *UserRepository) Add(user models.UserRequest) (*models.UserResponse, error) {
	// prepare the query to be used in the db.exec method
	const query = "INSERT INTO USERS (name, email, password, phone, is_adm) VALUES (?, ?, ?, ?, ?)"

	// execute the query prepared previously with the user columns
	res, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.IsAdm)

	// returns an error if any occurs
	if err != nil {
		return nil, err
	}

	// get the user id of the insert
	id, _ := res.LastInsertId()

	// returns a &models.UserResponse with the ID of the inserted user
	return &models.UserResponse{ID: id}, nil
}

// Delete a user from the database with the given userId. Return an error if any occurs
func (r *UserRepository) Delete(userId int64) error {
	// prepare the query to be used in the db.exec method
	const query = "DELETE FROM USERS WHERE ID = ?"

	// execute the query prepared previously with the user id
	_, err := r.DB.Exec(query, userId)

	// returns an error if any occurs
	if err != nil {
		return err
	}

	// function completed with no errors; return nil
	return nil
}

// FindById returns the user with the given userId in the database, or an error if any occurs
func (r *UserRepository) FindById(userId int64) (*models.User, error) {
	// prepare the query to be used in the db.QueryRow method
	const query = "SELECT id, name, email, password, phone, is_adm FROM USERS WHERE ID = ?"

	// execute the query prepared previously with the user id
	row := r.DB.QueryRow(query, userId)

	// instantiate a new user to be scanned
	var user models.User

	// try to scan the user from the given row
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.IsAdm); err != nil {

		// if no row is returned, then the given userId is not assigned to any user. Return err not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrNotFound
		}

		// return internal error
		return nil, core.ErrInternal
	}

	// return the scanned user
	return &user, nil
}
