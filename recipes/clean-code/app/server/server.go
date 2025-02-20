package server

import (
	"context"

	"app/datasources"
	"app/server/handlers"
	"app/server/services"

	"go.khulnasoft.com/velocity"
)

// NewServer creates a new Velocity app and sets up the routes
func NewServer(ctx context.Context, dataSources *datasources.DataSources) *velocity.App {
	app := velocity.New()
	apiRoutes := app.Group("/api")

	apiRoutes.Get("/status", func(c *velocity.Ctx) error {
		return c.SendString("ok")
	})
	apiRoutes.Get("/v1/books", handlers.GetBooks(services.NewBooksService(dataSources.DB)))
	apiRoutes.Post("/v1/books", handlers.AddBook(services.NewBooksService(dataSources.DB)))

	return app
}
