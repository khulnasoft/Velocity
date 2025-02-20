package router

import (
	"velocity-oauth-google/handler"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

// Routes for velocity
func Routes(app *velocity.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Auth)
	api.Get("/auth/google/callback", handler.Callback)
}
