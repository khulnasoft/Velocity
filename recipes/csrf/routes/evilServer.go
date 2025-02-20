package routes

import (
	"go.khulnasoft.com/velocity"
)

// RegisterEvilRoutes registers the routes and middlewares necessary for the server
func RegisterEvilRoutes(evilApp *fiber.App) {
	evilApp.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/evil-examples", fiber.Map{})
	})

	evilApp.Get("/malicious-form", func(c *fiber.Ctx) error {
		return c.Render("views/malicious-form", fiber.Map{})
	})
}
