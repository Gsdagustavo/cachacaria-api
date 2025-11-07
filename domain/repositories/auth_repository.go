package repositories

import (
	"cachacariaapi/domain/entities"
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	AddUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*entities.User, error)
	DeleteUser(ctx context.Context, id int) error
}
