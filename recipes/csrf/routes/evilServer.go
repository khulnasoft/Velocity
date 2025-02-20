package routes

import (
	"go.khulnasoft.com/velocity"
)

// RegisterEvilRoutes registers the routes and middlewares necessary for the server
func RegisterEvilRoutes(evilApp *velocity.App) {
	evilApp.Get("/", func(c *velocity.Ctx) error {
		return c.Render("views/evil-examples", velocity.Map{})
	})

	evilApp.Get("/malicious-form", func(c *velocity.Ctx) error {
		return c.Render("views/malicious-form", velocity.Map{})
	})
}
