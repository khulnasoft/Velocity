package handlers

import "go.khulnasoft.com/velocity"

// NotFound returns status code 404 along with the given html file
func NotFound(file string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendFile(file)
	}
}
