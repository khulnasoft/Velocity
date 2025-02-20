package handlers

import "go.khulnasoft.com/velocity"

// Render will pass the remove IP value to the template input
func Render() velocity.Handler {
	return func(c *velocity.Ctx) error {
		return c.Render("index", velocity.Map{
			"IP": c.IP(),
		})
	}
}
