package routes

import (
	"clean-architecture/api/handlers"
	"clean-architecture/pkg/book"

	"github.com/khulnasoft/fiber/v2"
)

// BookRouter is the Router for KhulnaSoft App
func BookRouter(app fiber.Router, service book.Service) {
	app.Get("/books", handlers.GetBooks(service))
	app.Post("/books", handlers.AddBook(service))
	app.Put("/books", handlers.UpdateBook(service))
	app.Delete("/books", handlers.RemoveBook(service))
}
