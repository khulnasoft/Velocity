package middleware

import (
	"strings"

	"numtostr/gotodo/utils/jwt"

	"go.khulnasoft.com/velocity"
)

// Auth is the authentication middleware
func Auth(c *velocity.Ctx) error {
	h := c.Get("Authorization")

	if h == "" {
		return velocity.ErrUnauthorized
	}

	// Spliting the header
	chunks := strings.Split(h, " ")

	// If header signature is not like `Bearer <token>`, then throw
	// This is also required, otherwise chunks[1] will throw out of bound error
	if len(chunks) < 2 {
		return velocity.ErrUnauthorized
	}

	// Verify the token which is in the chunks
	user, err := jwt.Verify(chunks[1])
	if err != nil {
		return velocity.ErrUnauthorized
	}

	c.Locals("USER", user.ID)

	return c.Next()
}
