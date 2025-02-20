package user

import (
	"context"

	"go.khulnasoft.com/velocity"
)

// If user does not exist, do not allow one to access the API.
func (h *UserHandler) checkIfUserExistsMiddleware(c *velocity.Ctx) error {
	// Create a new customized context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedUserID, err := c.ParamsInt("userID")
	if err != nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": "Please specify a valid user ID!",
		})
	}

	// Check if user exists.
	searchedUser, err := h.userService.GetUser(customContext, targetedUserID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	if searchedUser == nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": "There is no user with this ID!",
		})
	}

	// Store in locals for further processing in the real handler.
	c.Locals("userID", targetedUserID)
	return c.Next()
}
