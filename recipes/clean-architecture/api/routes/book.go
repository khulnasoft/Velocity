package routes

import (
	"clean-architecture/api/handlers"
	"clean-architecture/pkg/book"

	"go.khulnasoft.com/velocity"
)

// BookRouter is the Router for KhulnaSoft App
func BookRouter(app velocity.Router, service book.Service) {
	app.Get("/books", handlers.GetBooks(service))
	app.Post("/books", handlers.AddBook(service))
	app.Put("/books", handlers.UpdateBook(service))
	app.Delete("/books", handlers.RemoveBook(service))
}
