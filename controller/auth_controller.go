package controller

import (
	"auth-service/config"
	"auth-service/middleware"
	"auth-service/model"
	"auth-service/service"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	AuthService service.AuthService
	Validate    config.Validation
}

func NewAuthController(AuthService service.AuthService, validate config.Validation) AuthController {
	return AuthController{
		AuthService: AuthService,
		Validate:    validate,
	}
}

func (controller *AuthController) Route(app *fiber.App) {
	app.Post("register", controller.Register)
	app.Post("login", controller.Login)
	app.Get("token", controller.RefreshToken)
	app.Use(middleware.Verify())
	app.Post("logout", controller.Logout)
	app.Post("user", controller.GetAuthUser)
}

func (controller *AuthController) Register(ctx *fiber.Ctx) error {
	var request model.RegisterRequest
	ctx.BodyParser(&request)
	valid := controller.Validate
	errValidation := valid.ValidateRequest(request)
	if errValidation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GetResponse(fiber.StatusBadRequest, errValidation, "Validation Error"))
	}

	controller.AuthService.Register(request)
	return ctx.Status(fiber.StatusOK).JSON(nil)
}

func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var request model.LoginRequest
	ctx.BodyParser(&request)
	valid := controller.Validate
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
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Expires:  time.Now().Add(7200 * time.Hour),
	})
	response := model.GetResponse(responseCode, model.LoginResponse{
		AccessToken:       accessToken,
		RefreshTokenToken: refreshToken,
	}, "Success")

	return ctx.Status(responseCode).JSON(response)
}

func (controller *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	// tokenString := ctx.Get("Refresh-Token")
	tokenString := ""
	if ctx.Get("Refresh-Token") == "" {
		tokenString = ctx.Cookies("refresh_token")
	} else {
		tokenString = ctx.Get("Refresh-Token")
	}
	responseCode, accessToken := controller.AuthService.RefreshToken(tokenString)

	fmt.Println("token", ctx.Cookies("refresh_token"))
	if tokenString == "" {
		response := model.GetResponse(responseCode, nil, "UNAUTHORIZED")
		return ctx.Status(responseCode).JSON(response)
	}
	response := model.GetResponse(responseCode, model.LoginResponse{AccessToken: *accessToken, RefreshTokenToken: tokenString}, "success")
	return ctx.Status(responseCode).JSON(response)
}

func (controller *AuthController) Logout(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	responseCode := controller.AuthService.Logout(email)
	response := model.GetResponse(responseCode, nil, "success logout")
	return ctx.Status(responseCode).JSON(response)
}

func (controller *AuthController) GetAuthUser(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	responseCode, data := controller.AuthService.GetAuthUser(email)
	response := model.GetResponse(responseCode, data, "")
	return ctx.Status(responseCode).JSON(response)
}
