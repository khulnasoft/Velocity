package main

import (
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Use an external setup function in order
	// to configure the app in tests as well
	app := Setup()

	// start the application on http://localhost:3000
	log.Fatal(app.Listen(":3000"))
}

// Setup Setup a velocity app with all of its routes
func Setup() *velocity.App {
	// Initialize a new app
	app := velocity.New()

	// Register the index route with a simple
	// "OK" response. It should return status
	// code 200
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("OK")
	})

	// Return the configured app
	return app
}
