package service

import (
	"auth-service/entity"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/security"
	"fmt"
)

type AuthServiceImpl struct {
	UserRepo     repository.UserRepository
	RedisService RedisService
}

func NewAuthServiceImpl(userRepo repository.UserRepository, redisService RedisService) AuthService {
	return &AuthServiceImpl{
		UserRepo:     userRepo,
		RedisService: redisService,
	}
}

// Register implements AuthService.
func (service *AuthServiceImpl) Register(request model.RegisterRequest) {
	user := &entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     request.Email,
		Password:  security.GetHash([]byte(request.Password)),
	}

	err := service.UserRepo.CreateUser(user)
	if err != nil {

	}
}

// Login implements AuthService.
func (service *AuthServiceImpl) Login(request *model.LoginRequest) (responseCode int, accessToken, refreshToken string) {
	user := service.UserRepo.FindUserByEmail(request.Email)
	err := security.ComparePassword(user.Password, request.Password)
	if err != nil {
		return 400, "", ""
	}
	accessToken, _ = security.ClaimAccessToken(request.Email)
	refreshToken, _ = security.ClaimRefreshToken(request.Email)
	service.RedisService.Set("refresh_token:"+refreshToken, refreshToken)
	return 200, accessToken, refreshToken
}

func (service *AuthServiceImpl) RefreshToken(token string) (int, *string) {
	if token == "" {
		return 401, nil
	}

	byteToken := service.RedisService.Get("refresh_token:" + token)
	// redisRefreshToken := string(*byteToken)
	fmt.Println("byte token", len(*byteToken))
	if len(*byteToken) == 0 {
		return 403, nil
	}

	isValidToken := security.VerifyToken(token)
	if !isValidToken {
		return 401, nil
	}

	claims := security.DecodeTokenString(token)
	if claims == nil {
		return 403, nil
	}
	accessToken, err := security.ClaimAccessToken(claims["email"].(string))
	if err != nil {
		return 403, nil
	}
	return 200, &accessToken
}
