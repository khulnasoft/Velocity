package handlers

import (
	"fmt"

	"go.khulnasoft.com/velocity"
)

// Errors will process all errors returned to velocity
func Errors(file string) velocity.ErrorHandler {
	return func(c *velocity.Ctx, err error) error {
		fmt.Println(err.Error())
		return c.Status(velocity.StatusInternalServerError).SendFile(file)
	}
}
