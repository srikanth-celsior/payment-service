package main

import (
	"payment-service/handlers"

	"github.com/kataras/iris/v12"
)

func SetupRouter(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		ctx.WriteString("Payment Service is running")
	})
	app.Get("/healthz", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{"status": "ok"})
	})
	app.Get("/payments", handlers.GetAllPayments)
}
