package repositories

import (
	"cachacariaapi/domain/entities"
)

type UserRepository interface {
	GetAll() ([]entities.User, error)
	Add(user entities.User) error
	Delete(userId int64) error
	FindByEmail(email string) (*entities.User, error)
	FindById(userid int64) (*entities.User, error)
	Update(user entities.User, userId int64) error
}
