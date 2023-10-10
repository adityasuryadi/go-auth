package repository

import "auth-service/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	FindUserByEmail(email string) *entity.User
}
