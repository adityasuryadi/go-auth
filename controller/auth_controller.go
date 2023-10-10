package controller

import (
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
	responseCode, accessToken, refreshToken := controller.AuthService.Login(&request)
	return ctx.Status(responseCode).JSON(model.LoginResponse{AccessToken: accessToken, RefreshTokenToken: refreshToken})
}

// type RefreshTokenRequest struct {
// 	RefreshToken string `json:"refresh_token"`
// }

func (controller *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	// refreshToken := ctx.Cookies("refresh_token")
	// fmt.Println(refreshToken)
	// cookie := new(fiber.Cookie)
	// cookie.Name = "refresh_token"
	// cookie.Value = "wkwkwkkw"
	// cookie.HTTPOnly = true
	// cookie.Expires = time.Now().Add(24 * time.Hour)
	// ctx.Cookie(cookie)
	// return ctx.Status(200).JSON("wkwkwkkw")
	// redis := config.NewRedis()
	// token, err := redis.Get("refresh_token:c45469ae-dd47-43f2-bf79-cc641323d426")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return ctx.SendString(string(token))

	tokenString := ctx.Get("Refresh-Token")
	responseCode, accessToken := controller.AuthService.RefreshToken(tokenString)
	// byteToken := service.RedisService.Get("refresh_token:"+token)
	return ctx.Status(responseCode).JSON(accessToken)
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	// Don't forget to validate the alg is what you expect:
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}

	// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	// 	return []byte("rahasia"), nil
	// })

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	// 	fmt.Println(claims)
	// } else {
	// 	fmt.Println(err)
	// }

	return nil
}
