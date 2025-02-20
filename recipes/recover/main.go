// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"log"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/recover"
)

func main() {
	// Velocity instance
	app := velocity.New(velocity.Config{
		// ErrorHandler: func(c *velocity.Ctx, err error) error {
		// 	c.Status(velocity.StatusInternalServerError)
		// 	return c.SendString(err.Error())
		// },
	})

	// Middleware
	app.Use(recover.New())

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler
func hello(c *velocity.Ctx) error {
	panic("No worries, I won't crash! 🙏")
}
