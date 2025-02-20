package utils

import (
	"go.khulnasoft.com/velocity"
)

// ParseBody is helper function for parsing the body.
// Is any error occurs it will panic.
// Its just a helper function to avoid writing if condition again n again.
func ParseBody(ctx *velocity.Ctx, body interface{}) *velocity.Error {
	if err := ctx.BodyParser(body); err != nil {
		return velocity.ErrBadRequest
	}

	return nil
}

// ParseBodyAndValidate is helper function for parsing the body.
// Is any error occurs it will panic.
// Its just a helper function to avoid writing if condition again n again.
func ParseBodyAndValidate(ctx *velocity.Ctx, body interface{}) *velocity.Error {
	if err := ParseBody(ctx, body); err != nil {
		return err
	}

	return Validate(body)
}

// GetUser is helper function for getting authenticated user's id
func GetUser(c *velocity.Ctx) *uint {
	id, _ := c.Locals("USER").(uint)
	return &id
}
