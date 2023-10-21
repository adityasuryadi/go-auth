package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/middleware"
	"auth-service/repository"
	"auth-service/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Panic occured", r)
	// 	} else {
	// 		fmt.Println("Application running perfectly")
	// 	}
	// }()

	app := fiber.New()
	configuration := config.New(".env")
	db := config.NewPostgresDB(configuration)
	redisConfig := config.NewRedis(configuration)
	userRepo := repository.NewUserRepository(db)
	redisService := service.NewRedisConfig(redisConfig, 7200)
	authService := service.NewAuthService(userRepo, redisService)
	authController := controller.NewAuthController(authService)
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Range, Authorization",
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))
	authController.Route(app)
	app.Use(middleware.Verify())
	app.Get("/test", func(c *fiber.Ctx) error {
		fmt.Println("user ", c.Locals("user"))
		return c.SendString("AUth, Service")
	})

	// test := new(Mamalia)
	// test.Route(app)
	app.Listen(":5001")
}
