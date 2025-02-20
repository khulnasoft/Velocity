package router

import (
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/limiter"
)

type ApiRouter struct{}

func (h ApiRouter) InstallRouter(app *velocity.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *velocity.Ctx) error {
		return ctx.Status(velocity.StatusOK).JSON(velocity.Map{
			"message": "Hello from api",
		})
	})
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
