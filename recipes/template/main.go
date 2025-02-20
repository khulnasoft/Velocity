package main

import (
	"log"

	"github.com/khulnasoft/template/django/v3"
	"go.khulnasoft.com/velocity"
)

func main() {
	// Create a new engine
	engine := django.New("./views", ".html")

	// Or from an embedded system
	// See github.com/khulnasoft/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".django"))

	// Pass the engine to the Views
	app := velocity.New(velocity.Config{
		Views: engine,
	})

	app.Get("/", func(c *velocity.Ctx) error {
		// Render with and extends
		return c.Render("index", velocity.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/embed", func(c *velocity.Ctx) error {
		// Render index within layouts/main
		return c.Render("embed", velocity.Map{
			"Title": "Hello, World!",
		}, "layouts/main2")
	})

	log.Fatal(app.Listen(":3000"))
}
