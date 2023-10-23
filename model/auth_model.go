package model

type RegisterRequest struct {
	FirstName            string `json:"first_name" validate:"required,alpha"`
	LastName             string `json:"last_name" validate:"required,alpha"`
	Email                string `json:"email" validate:"required,email,unique=users"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=password"`
	Phone                string `json:"phone" validate:"required,numeric"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken       string `json:"access_token"`
	RefreshTokenToken string `json:"refresh_token"`
}

type UserResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}
