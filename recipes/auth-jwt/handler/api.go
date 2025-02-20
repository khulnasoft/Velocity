package handler

import "go.khulnasoft.com/velocity"

// Hello handle api status
func Hello(c *velocity.Ctx) error {
	return c.JSON(velocity.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
