package auth

import (
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	jwtware "github.com/khulnasoft/contrib/jwt"
	"go.khulnasoft.com/velocity"
)

// JWT error message.
func jwtError(c *velocity.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(velocity.StatusBadRequest).JSON(&velocity.Map{
			"status":  "error",
			"message": "Missing or malformed JWT!",
		})
	}

	return c.Status(velocity.StatusUnauthorized).JSON(&velocity.Map{
		"status":  "error",
		"message": "Invalid or expired JWT!",
	})
}

// Guards a specific endpoint in the API.
func JWTMiddleware() velocity.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: jwtError,
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		TokenLookup:  "cookie:jwt",
	})
}

// Gets user data (their ID) from the JWT middleware. Should be executed after calling 'JWTMiddleware()'.
func GetDataFromJWT(c *velocity.Ctx) error {
	// Get userID from the previous route.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)
	parsedUserID := claims["uid"].(string)
	userID, err := strconv.Atoi(parsedUserID)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Go to next.
	c.Locals("currentUser", userID)
	return c.Next()
}
