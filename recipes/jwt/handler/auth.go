package handler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.khulnasoft.com/velocity"
)

// Login get user and password
func Login(c *velocity.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(velocity.StatusUnauthorized)
	}
	identity := input.Identity
	pass := input.Password
	if identity != "ender" || pass != "ender" {
		return c.SendStatus(velocity.StatusUnauthorized)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = identity
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(velocity.StatusInternalServerError)
	}

	return c.JSON(velocity.Map{"status": "success", "message": "Success login", "data": t})
}
