package handlers

import "go.khulnasoft.com/velocity"

// NotFound returns status code 404 along with the given html file
func NotFound(file string) velocity.Handler {
	return func(c *velocity.Ctx) error {
		return c.Status(velocity.StatusNotFound).SendFile(file)
	}
}
