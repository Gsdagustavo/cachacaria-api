package repositories

import "cachacariaapi/internal/domain/entities"

type UserRepository interface {
	GetAll() ([]entities.User, error)
	Add(user entities.RegisterRequest) (*entities.UserResponse, error)
	Delete(userId int64) error
	FindByEmail(email string) (*entities.User, error)
	FindById(userid int64) (*entities.User, error)
	Update(user entities.UserRequest, userId int64) (*entities.UserResponse, error)
}
