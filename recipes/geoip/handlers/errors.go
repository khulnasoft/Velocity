package handlers

import (
	"fmt"

	"github.com/khulnasoft/fiber/v2"
)

// Errors will process all errors returned to fiber
func Errors(file string) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).SendFile(file)
	}
}
