package src

import (
	"example.com/KhulnasoftFirebaseBoilerplate/src/routes"

	"go.khulnasoft.com/velocity"
)

func CreateServer() *velocity.App {
	version := "v1.0.0"

	app := velocity.New(velocity.Config{
		ServerHeader: "Khulnasoft Firebase Boilerplate",
		AppName:      "Khulnasoft Firebase Boilerplate " + version,
	})

	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Khulnasoft Firebase Boilerplate " + version)
	})

	routes.New().Setup(app)

	return app
}
