package handlers

import (
	"go.khulnasoft.com/velocity"
)

// Home renders the home view
func Home(c *velocity.Ctx) error {
	return c.Render("index", velocity.Map{
		"Title": "Hello, World!",
	})
}

// About renders the about view
func About(c *velocity.Ctx) error {
	return c.Render("about", nil)
}

// NoutFound renders the 404 view
func NotFound(c *velocity.Ctx) error {
	return c.Status(404).Render("404", nil)
}
