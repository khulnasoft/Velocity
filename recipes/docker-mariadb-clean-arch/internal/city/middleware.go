package city

import (
	"context"

	"go.khulnasoft.com/velocity"
)

// If city does not exist, do not allow one to access the API.
func (h *CityHandler) checkIfCityExistsMiddleware(c *velocity.Ctx) error {
	// Create a new customized context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedCityID, err := c.ParamsInt("cityID")
	if err != nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": "Please specify a valid city ID!",
		})
	}

	// Check if city exists.
	searchedCity, err := h.cityService.FetchCity(customContext, targetedCityID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	if searchedCity == nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": "There is no city with this ID!",
		})
	}

	// Store in locals for further processing in the real handler.
	c.Locals("cityID", targetedCityID)
	return c.Next()
}
