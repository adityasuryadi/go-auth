package test

import (
	"auth-service/config"
	"auth-service/entity"
	"auth-service/mock"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userRepositoryMock = &mock.UserRepositoryMock{}

func TestRegisterFail(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	request := model.RegisterRequest{
		FirstName: "",
		LastName:  "",
		Phone:     "",
		Email:     "",
		Password:  "",
	}
	authService := service.NewAuthService(userRepo, redisService)
	responseCode := authService.Register(request)
	assert.Equal(t, responseCode, 500)
}

func TestRegisterSuccess(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)
	user := entity.User{
		FirstName: "Aditya",
		LastName:  "Suryadi",
		Phone:     "081214124",
		Email:     "adit@mail.com",
		Password:  "123456",
	}

	request := model.RegisterRequest{
		FirstName: "Aditya",
		LastName:  "Suryadi",
		Phone:     "081214124",
		Email:     "adit@mail.com",
		Password:  "123456",
	}

	userRepositoryMock.Mock.On("CreateUser", user).Return(nil)
	authService := service.NewAuthService(userRepo, redisService)
	responseCode := authService.Register(request)
	assert.Equal(t, responseCode, 200)
}

func TestLoginUserNotfound(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	request := &model.LoginRequest{
		Email:    "notfound@mail.com",
		Password: "password",
	}
	userService := service.NewAuthService(userRepo, redisService)
	responseCode, accessToken, refreshToken := userService.Login(request)
	assert.Equal(t, responseCode, 400)
	assert.Equal(t, accessToken, "")
	assert.Equal(t, refreshToken, "")
}

func TestLoginUserWrongPassword(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	request := &model.LoginRequest{
		Email:    "adit@mail.com",
		Password: "passwordsss",
	}
	userService := service.NewAuthService(userRepo, redisService)
	responseCode, accessToken, refreshToken := userService.Login(request)
	assert.Equal(t, responseCode, 400)
	assert.Equal(t, accessToken, "")
	assert.Equal(t, refreshToken, "")
}

func TestLoginUserSuccess(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	request := &model.LoginRequest{
		Email:    "adit@mail.com",
		Password: "password",
	}
	userService := service.NewAuthService(userRepo, redisService)
	responseCode, accessToken, refreshToken := userService.Login(request)
	assert.Equal(t, responseCode, 200)
	assert.NotEqual(t, accessToken, "")
	assert.NotEqual(t, refreshToken, "")
}

func TestRefreshTokenWithoutToken(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	userService := service.NewAuthService(userRepo, redisService)
	responseCode, accessToken := userService.RefreshToken("")
	assert.Equal(t, responseCode, 401)
	assert.Nil(t, accessToken)
}

func TestRefreshTokenWithInvalidToken(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	userService := service.NewAuthService(userRepo, redisService)
	responseCode, accessToken := userService.RefreshToken("skjjfqflasjdiasdaidhasdjad")
	assert.Equal(t, responseCode, 403)
	assert.Nil(t, accessToken)
}

func TestRefreshTokenWithValidToken(t *testing.T) {
	configuration := config.New(`.\..\.env.test`)
	db := config.NewPostgresDB(configuration)
	userRepo := repository.NewUserRepository(db)
	redisConfig := config.NewRedis(configuration)
	redisService := service.NewRedisConfig(redisConfig, 7200)

	userService := service.NewAuthService(userRepo, redisService)
	request := &model.LoginRequest{
		Email:    "adit@mail.com",
		Password: "password",
	}
	_, _, refreshToken := userService.Login(request)
	responseCode, accessToken := userService.RefreshToken(refreshToken)
	assert.Equal(t, responseCode, 200)
	assert.NotEqual(t, accessToken, "")
}
