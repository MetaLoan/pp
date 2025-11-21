package repository

import (
	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return core.DB.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := core.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserByUUID(uuid string) (*model.User, error) {
	var user model.User
	err := core.DB.Where("uuid = ?", uuid).First(&user).Error
	return &user, err
}
