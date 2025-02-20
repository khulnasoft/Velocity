package route

import (
	"velocity-sqlboiler/api/controller"

	"go.khulnasoft.com/velocity"
)

func SetupRoutes(app *velocity.App) {
	api := app.Group("/api/v1")
	// Post
	api.Get("/posts", controller.GetPosts)
	api.Get("/posts/:id", controller.GetPost)
	api.Post("/posts", controller.NewPost)
	api.Delete("/posts/:id", controller.DeletePost)
	api.Put("/posts/:id", controller.UpdatePost)

	// Author
	api.Get("/authors", controller.GetAuthors)
	api.Get("/authors/:id", controller.GetAuthor)
	api.Post("/authors", controller.NewAuthor)
	api.Delete("/authors/:id", controller.DeleteAuthor)
	api.Put("/authors/:id", controller.UpdateAuthor)
}
