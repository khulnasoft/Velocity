package main

import (
	"log"

	"go.khulnasoft.com/velocity"                    // Importing the velocity package for handling HTTP requests
	"go.khulnasoft.com/velocity/middleware/cors"    // Middleware for handling Cross-Origin Resource Sharing (CORS)
	"go.khulnasoft.com/velocity/middleware/favicon" // Middleware for serving favicon
	"go.khulnasoft.com/velocity/middleware/logger"  // Middleware for logging HTTP requests
	// Package for logging errors
)

func main() {
	app := velocity.New() // Initialize a new Velocity instance
	// register middlewares
	app.Use(favicon.New()) // Use favicon middleware to serve favicon
	app.Use(cors.New())    // Use CORS middleware to allow cross-origin requests
	app.Use(logger.New())  // Use logger middleware to log HTTP requests

	// Define a GET route for the path '/hello'
	app.Get("/hello", func(c *velocity.Ctx) error {
		return c.SendString("World!") // Send a response when the route is accessed
	})

	log.Fatal(app.Listen(":5000")) // Start the server on port 5000 and log any errors
}
