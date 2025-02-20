package middleware

import (
	jwtware "github.com/khulnasoft/contrib/jwt"
	"go.khulnasoft.com/velocity"
)

// Protected protect routes
func Protected() func(*velocity.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte("secret")},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *velocity.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(velocity.StatusBadRequest)
		return c.JSON(velocity.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(velocity.StatusUnauthorized)
		return c.JSON(velocity.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
