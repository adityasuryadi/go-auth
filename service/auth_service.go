package service

import "auth-service/model"

type AuthService interface {
	Register(request model.RegisterRequest)
	RefreshToken(token string) (int, *string)
	Login(request *model.LoginRequest) (responseCode int, accessToken, refreshToken string)
}
