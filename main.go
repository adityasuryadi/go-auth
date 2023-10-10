package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/middleware"
	"auth-service/repository"
	"auth-service/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// type animal interface {
// 	Sound() string
// }

// type Mamalia struct {
// 	animal animals.Duck
// }

// func (a Mamalia) SuaraAnjing(ctx *fiber.Ctx) error {
// 	suara := a.animal.Sound()
// 	return ctx.SendString(suara)
// }

//	func (a Mamalia) Play(c *fiber.Ctx) error {
//		sound := a.animal.Sound()
//		return c.
//	}
// func (a Mamalia) Route(app *fiber.App) {
// 	app.Get("suara", a.SuaraAnjing)
// }

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occured", r)
		} else {
			fmt.Println("Application running perfectly")
		}
	}()

	app := fiber.New()
	configuration := config.New(".env")
	db := config.NewPostgresDB(configuration)
	redisConfig := config.NewRedis()
	userRepo := repository.NewUserRepository(db)
	redisService := service.NewRedisConfig(redisConfig, 7200)
	authService := service.NewAuthServiceImpl(userRepo, redisService)
	authController := controller.NewAuthController(authService)
	authController.Route(app)
	app.Use(middleware.Verify())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("AUth, Service")
	})

	// test := new(Mamalia)
	// test.Route(app)
	app.Listen(":5001")
}
