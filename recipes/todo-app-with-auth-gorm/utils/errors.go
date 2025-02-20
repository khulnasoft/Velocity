package utils

import (
	"go.khulnasoft.com/velocity"
)

type httpError struct {
	Statuscode int    `json:"statusCode"`
	Error      string `json:"error"`
}

// ErrorHandler is used to catch error thrown inside the routes by ctx.Next(err)
func ErrorHandler(c *velocity.Ctx, err error) error {
	// Statuscode defaults to 500
	code := velocity.StatusInternalServerError

	// Check if it's an velocity.Error type
	if e, ok := err.(*velocity.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(&httpError{
		Statuscode: code,
		Error:      err.Error(),
	})
}
