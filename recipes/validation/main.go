package main

import (
	"log"

	"validation/config"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
)

func main() {
	app := velocity.New()
	validate := validator.New() // Create Validate for using.

	// Use Cors
	app.Use(cors.New())

	app.Get("/test", func(ctx *velocity.Ctx) error {
		type User struct {
			ID        uint   `validate:"required,omitempty"`
			Firstname string `validate:"required"`
			Password  string `validate:"gte=10"` // gte = Greater than or equal
		}

		user := User{
			ID:        1,
			Firstname: "Velocity",
			/*
				if you delete Firstname field
				you'll get response like this: Error:Field validation for 'Firstname' failed on the 'required' tag"
			*/
			Password: "VelocityPassword123",
			/*
				if you enter "Velocity" in Password
				you'll get response like this: Error:Field validation for 'Password' failed on the 'gte' tag"
			*/
		}

		if err := validate.Struct(user); err != nil {
			return ctx.Status(velocity.StatusBadRequest).JSON(err.Error())
		}

		return ctx.Status(velocity.StatusOK).JSON("success time")
	})

	// .env Variables validation
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	log.Fatal(app.Listen(config.Config("PORT")))
}
