package src

import (
	"example.com/KhulnasoftFirebaseBoilerplate/src/routes"

	"go.khulnasoft.com/velocity"
)

func CreateServer() *fiber.App {
	version := "v1.0.0"

	app := fiber.New(fiber.Config{
		ServerHeader: "Khulnasoft Firebase Boilerplate",
		AppName:      "Khulnasoft Firebase Boilerplate " + version,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Khulnasoft Firebase Boilerplate " + version)
	})

	routes.New().Setup(app)

	return app
}
