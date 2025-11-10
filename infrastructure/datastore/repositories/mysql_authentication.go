package repositories

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/infrastructure/datastore"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type MySQLAuthRepository struct {
	db *sql.DB
}

func NewMySQLAuthRepository(db *sql.DB) repositories.AuthRepository {
	return &MySQLAuthRepository{
		db: db,
	}
}

func (r MySQLAuthRepository) AddUser(ctx context.Context, user *entities.User) error {
	const query = `
		INSERT INTO users (uuid, email, password, phone, is_adm) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.UUID,
		user.Email,
		user.Password,
		user.Phone,
		user.IsAdm,
	)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to add user"), err)
	}

	return nil
}

func (r MySQLAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	const query = `
		SELECT id, uuid, email, password, phone, is_adm FROM users WHERE email = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.UUID, &user.Email, &user.Password, &user.Phone, &user.IsAdm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query/scan user"), err)
	}

	return &user, nil
}

func (r MySQLAuthRepository) GetUserByPhone(
	ctx context.Context,
	phone string,
) (*entities.User, error) {
	const query = `
		SELECT id, uuid, email, password, phone, is_adm FROM users WHERE phone = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, phone).
		Scan(&user.ID, &user.UUID, &user.Email, &user.Password, &user.Phone, &user.IsAdm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query/scan user"), err)
	}

	return &user, nil
}

func (r MySQLAuthRepository) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	const query = `
		SELECT id, uuid, email, password, phone, is_adm FROM users WHERE id = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.UUID, &user.Email, &user.Password, &user.Phone, &user.IsAdm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query/scan user"), err)
	}

	return &user, nil
}

func (r MySQLAuthRepository) GetUserByUUID(
	ctx context.Context,
	uuid uuid.UUID,
) (*entities.User, error) {
	const query = `
		SELECT id, uuid, email, password, phone, is_adm FROM users WHERE uuid = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, uuid).
		Scan(&user.ID, &user.UUID, &user.Email, &user.Password, &user.Phone, &user.IsAdm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(fmt.Errorf("failed to query/scan user"), err)
	}

	return &user, nil
}

func (r MySQLAuthRepository) DeleteUser(ctx context.Context, id int) error {
	const query = `
		DELETE FROM users WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to delete user"), err)
	}

	return nil
}
