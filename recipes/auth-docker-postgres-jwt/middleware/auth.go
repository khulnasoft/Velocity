package middleware

import (
	"app/config"

	jwtware "github.com/khulnasoft/contrib/jwt"
	"go.khulnasoft.com/velocity"
)

// Protected protect routes
func Protected() velocity.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *velocity.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(velocity.StatusBadRequest).
			JSON(velocity.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(velocity.StatusUnauthorized).
		JSON(velocity.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
