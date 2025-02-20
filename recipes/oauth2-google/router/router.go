package router

import (
	"fiber-oauth-google/handler"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

// Routes for fiber
func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Auth)
	api.Get("/auth/google/callback", handler.Callback)
}
