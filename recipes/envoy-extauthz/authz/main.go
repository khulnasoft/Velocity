package main

import (
	"log"
	"strings"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/keyauth"
	"go.khulnasoft.com/velocity/middleware/logger"
)

const (
	authKey  = "apiKey"
	authName = "x-api-key"
	authSrc  = "header"
)

var (
	authList   = []string{"valid-key"}
	errMissing = &velocity.Error{
		Code:    403000,
		Message: "Missing API key",
	}
	errInvalid = &velocity.Error{
		Code:    403001,
		Message: "Invalid API key",
	}
)

func main() {
	app := velocity.New()

	app.Use(logger.New())

	app.Use(keyauth.New(keyauth.Config{
		SuccessHandler: successHandler,
		ErrorHandler:   errHandler,
		KeyLookup:      strings.Join([]string{authSrc, authName}, ":"),
		Validator:      validator,
		ContextKey:     authKey,
	}))

	log.Fatal(app.Listen(":1337"))
}

func successHandler(ctx *velocity.Ctx) error {
	return ctx.SendStatus(velocity.StatusOK)
}

func errHandler(ctx *velocity.Ctx, err error) error {
	ctx.Status(velocity.StatusForbidden)

	if err == errMissing {
		return ctx.JSON(errMissing)
	}

	return ctx.JSON(errInvalid)
}

func validator(ctx *velocity.Ctx, s string) (bool, error) {
	if s == "" {
		return false, errMissing
	}

	for _, val := range authList {
		if s == val {
			return true, nil
		}
	}

	return false, errInvalid
}
