package service

import "auth-service/model"

type AuthService interface {
	Register(request model.RegisterRequest) (responseCode int)
	RefreshToken(token string) (int, *string)
	Login(request *model.LoginRequest) (responseCode int, accessToken, refreshToken string)
	Logout(email string) (responseCode int)
	GetAuthUser(email string) (responseCode int, response model.UserResponse)
}
