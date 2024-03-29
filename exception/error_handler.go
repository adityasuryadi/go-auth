package exception

import (
	// "payment/config"

	"auth-service/model"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/natefinch/lumberjack"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./var/log/application.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	})
	log.Print(err.Error())

	// _, ok := err.(ValidationError)
	// if ok {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse{
	// 		Code:   400,
	// 		Status: "BAD_REQUEST",
	// 		Data:   err.Error(),
	// 	})
	// }
	return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse{
		Code:   500,
		Status: "INTERNAL_SERVER_ERROR",
		Data:   err.Error(),
	})
}

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}
}
