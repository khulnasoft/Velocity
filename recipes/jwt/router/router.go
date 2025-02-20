package router

import (
	"api-velocity-gorm/handler"
	"api-velocity-gorm/middleware"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *velocity.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// Products
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)
}
