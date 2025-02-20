package main

import (
	"log"
	"os"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Initialize the application
	app := velocity.New()

	// Hello, World!
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Listen and Server in 0.0.0.0:$PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen(":" + port))
}
