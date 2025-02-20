package router

import (
	"fiber-oauth-google/handler"

	"github.com/khulnasoft/fiber/v2"
	"github.com/khulnasoft/fiber/v2/middleware/logger"
)

// Routes for fiber
func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Auth)
	api.Get("/auth/google/callback", handler.Callback)
}
