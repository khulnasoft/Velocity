package city

import (
	"context"

	"docker-mariadb-clean-arch/internal/auth"

	"go.khulnasoft.com/velocity"
)

// We will inject our dependency - the service - here.
type CityHandler struct {
	cityService CityService
}

// Creates a new handler.
func NewCityHandler(cityRoute velocity.Router, cs CityService) {
	// Create a handler based on our created service / use-case.
	handler := &CityHandler{
		cityService: cs,
	}

	// We will restrict this route with our JWT middleware.
	// You can inject other middlewares if you see fit here.
	cityRoute.Use(auth.JWTMiddleware(), auth.GetDataFromJWT)

	// Routing for general routes.
	cityRoute.Get("", handler.getCities)
	cityRoute.Post("", handler.createCity)

	// Routing for specific routes.
	cityRoute.Get("/:cityID", handler.getCity)
	cityRoute.Put("/:cityID", handler.checkIfCityExistsMiddleware, handler.updateCity)
	cityRoute.Delete("/:cityID", handler.checkIfCityExistsMiddleware, handler.deleteCity)
}

// Handler to get all cities.
func (h *CityHandler) getCities(c *velocity.Ctx) error {
	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get all cities.
	cities, err := h.cityService.FetchCities(customContext)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return results.
	return c.Status(velocity.StatusOK).JSON(&velocity.Map{
		"status": "success",
		"data":   cities,
	})
}

// Get one city.
func (h *CityHandler) getCity(c *velocity.Ctx) error {
	// Create cancellable context.
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

	// Get one city.
	city, err := h.cityService.FetchCity(customContext, targetedCityID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return results.
	return c.Status(velocity.StatusOK).JSON(&velocity.Map{
		"status": "success",
		"data":   city,
	})
}

// Creates a single city.
func (h *CityHandler) createCity(c *velocity.Ctx) error {
	// Initialize variables.
	city := &City{}
	currentUserID := c.Locals("currentUser").(int)

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse request body.
	if err := c.BodyParser(city); err != nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Create one city.
	err := h.cityService.BuildCity(customContext, city, currentUserID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return result.
	return c.Status(velocity.StatusCreated).JSON(&velocity.Map{
		"status":  "success",
		"message": "City has been created successfully!",
	})
}

// Updates a single city.
func (h *CityHandler) updateCity(c *velocity.Ctx) error {
	// Initialize variables.
	city := &City{}
	currentUserID := c.Locals("currentUser").(int)
	targetedCityID := c.Locals("cityID").(int)

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse request body.
	if err := c.BodyParser(city); err != nil {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Update one city.
	err := h.cityService.ModifyCity(customContext, targetedCityID, city, currentUserID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return result.
	return c.Status(velocity.StatusOK).JSON(&velocity.Map{
		"status":  "success",
		"message": "City has been updated successfully!",
	})
}

// Deletes a single city.
func (h *CityHandler) deleteCity(c *velocity.Ctx) error {
	// Initialize previous city ID.
	targetedCityID := c.Locals("cityID").(int)

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete one city.
	err := h.cityService.DestroyCity(customContext, targetedCityID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return 204 No Content.
	return c.SendStatus(velocity.StatusNoContent)
}
