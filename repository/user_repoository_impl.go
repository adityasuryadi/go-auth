package repository

import (
	"auth-service/entity"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: *db,
	}
}

type UserRepositoryImpl struct {
	db gorm.DB
}

// FindUserByEmail implements UserRepository.
func (repository *UserRepositoryImpl) FindUserByEmail(email string) *entity.User {
	var user *entity.User
	err := repository.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil
	}
	return user
}

// CreateUser implements UserRepository.
func (repository *UserRepositoryImpl) CreateUser(user *entity.User) error {
	err := repository.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
