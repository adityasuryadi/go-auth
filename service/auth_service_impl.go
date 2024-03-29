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

// GetAuthUser implements AuthService.
func (service *AuthServiceImpl) GetAuthUser(email string) (responseCode int, response model.UserResponse) {
	user := service.UserRepo.FindUserByEmail(email)
	if user != nil {
		response = model.UserResponse{
			Id:        user.Id.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		}
		return 200, response
	}
	return 404, response
}

// Logout implements AuthService.
func (service *AuthServiceImpl) Logout(email string) (responseCode int) {
	err := service.RedisService.Delete("refresh_token:" + email)
	if err != nil {
		return 500
	}
	return 200
}

func NewAuthService(userRepo repository.UserRepository, redisService RedisService) AuthService {
	return &AuthServiceImpl{
		UserRepo:     userRepo,
		RedisService: redisService,
	}
}

// Register implements AuthService.
func (service *AuthServiceImpl) Register(request model.RegisterRequest) (responseCode int) {
	user := &entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     request.Email,
		Password:  security.GetHash([]byte(request.Password)),
	}

	err := service.UserRepo.CreateUser(user)
	if err != nil {
		return 200
	}
	return 500
}

// Login implements AuthService.
func (service *AuthServiceImpl) Login(request *model.LoginRequest) (responseCode int, accessToken, refreshToken string) {
	user := service.UserRepo.FindUserByEmail(request.Email)
	if user == nil {
		return 400, "", ""
	}
	err := security.ComparePassword(user.Password, request.Password)
	if err != nil {
		return 400, "", ""
	}
	accessToken, _ = security.ClaimAccessToken(request.Email)
	refreshToken, _ = security.ClaimRefreshToken(request.Email)
	service.RedisService.Set("refresh_token:"+user.Email, refreshToken)
	return 200, accessToken, refreshToken
}

func (service *AuthServiceImpl) RefreshToken(token string) (int, *string) {
	if token == "" {
		return 401, nil
	}

	jwtToken, isValidToken := security.VerifyToken(token)
	if !isValidToken {
		return 401, nil
	}

	claims := security.DecodeToken(jwtToken)
	if claims == nil {
		return 403, nil
	}

	byteToken := service.RedisService.Get("refresh_token:" + claims["email"].(string))
	// redisRefreshToken := string(*byteToken)
	fmt.Println("byte token", len(*byteToken))
	if len(*byteToken) == 0 {
		return 403, nil
	}

	accessToken, err := security.ClaimAccessToken(claims["email"].(string))
	if err != nil {
		return 403, nil
	}
	return 200, &accessToken
}
