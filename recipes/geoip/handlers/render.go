package handlers

import "go.khulnasoft.com/velocity"

// Render will pass the remove IP value to the template input
func Render() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"IP": c.IP(),
		})
	}
}
