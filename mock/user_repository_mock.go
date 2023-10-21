package mock

import (
	"auth-service/entity"
	"errors"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) CreateUser(user *entity.User) error {
	arguments := repository.Mock.Called(user)
	if arguments.Get(0) == nil {
		return errors.New("failed create user")
	} else {
		return nil
	}
}

func (repository *UserRepositoryMock) FindUserByEmail(email string) *entity.User {
	arguments := repository.Mock.Called(email)
	if arguments.Get(0) == nil {
		return nil
	}
	var user *entity.User
	return user
}
