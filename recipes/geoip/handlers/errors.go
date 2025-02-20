package handlers

import (
	"fmt"

	"go.khulnasoft.com/velocity"
)

// Errors will process all errors returned to fiber
func Errors(file string) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).SendFile(file)
	}
}
