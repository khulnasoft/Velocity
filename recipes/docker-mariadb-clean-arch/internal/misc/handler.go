package misc

import "go.khulnasoft.com/velocity"

// Create a handler. Leave this empty, as we have no domains nor use-cases.
type MiscHandler struct{}

// Represents a new handler.
func NewMiscHandler(miscRoute velocity.Router) {
	handler := &MiscHandler{}

	// Declare routing.
	miscRoute.Get("", handler.healthCheck)
}

// Check for the health of the API.
func (h *MiscHandler) healthCheck(c *velocity.Ctx) error {
	return c.Status(velocity.StatusOK).JSON(&velocity.Map{
		"status":  "success",
		"message": "Hello World!",
	})
}
