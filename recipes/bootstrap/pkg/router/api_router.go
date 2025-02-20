package router

import (
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/limiter"
)

type ApiRouter struct{}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
