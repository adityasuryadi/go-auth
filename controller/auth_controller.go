package controller

import (
	"auth-service/config"
	"auth-service/model"
	"auth-service/service"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(AuthService service.AuthService) AuthController {
	return AuthController{
		AuthService: AuthService,
	}
}

func (controller *AuthController) Route(app *fiber.App) {
	app.Post("register", controller.Register)
	app.Post("login", controller.Login)
	app.Get("token", controller.RefreshToken)
}

func (controller *AuthController) Register(ctx *fiber.Ctx) error {
	var request model.RegisterRequest
	ctx.BodyParser(&request)
	controller.AuthService.Register(request)
	return ctx.Status(fiber.StatusOK).JSON(nil)
}

func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var request model.LoginRequest
	ctx.BodyParser(&request)
	valid := config.NewValidation()
	errValidation := valid.ValidateRequest(request)
	if errValidation != nil {
		return ctx.Status(400).JSON(model.GetResponse(400, errValidation, "Validation Error"))
	}
	responseCode, accessToken, refreshToken := controller.AuthService.Login(&request)
	if accessToken == "" {
		return ctx.Status(400).JSON(model.GetResponse(400, "", "Email Or Password Incorect"))
	}
	// var data = []interface{}{
	// 	model.LoginResponse{
	// 		AccessToken:       accessToken,
	// 		RefreshTokenToken: refreshToken,
	// 	},
	// }
	response := model.GetResponse(responseCode, model.LoginResponse{
		AccessToken:       accessToken,
		RefreshTokenToken: refreshToken,
	}, "Success")
	return ctx.Status(responseCode).JSON(response)
}

func (controller *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Refresh-Token")
	responseCode, accessToken := controller.AuthService.RefreshToken(tokenString)
	response := model.GetResponse(responseCode, model.LoginResponse{AccessToken: *accessToken, RefreshTokenToken: tokenString}, "success")
	return ctx.Status(responseCode).JSON(response)
}
