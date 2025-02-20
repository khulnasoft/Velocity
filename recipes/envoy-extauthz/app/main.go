package main

import (
	"log"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

func main() {
	app := velocity.New()

	app.Use(logger.New())

	app.Get("/health", func(ctx *velocity.Ctx) error {
		return ctx.SendString("Healthy")
	})

	api := app.Group("/api")

	api.Get("/resource", func(ctx *velocity.Ctx) error {
		return ctx.SendString("Some Resource API")
	})

	log.Fatal(app.Listen(":3000"))
}
