// 🚀 Velocity is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://docs.khulnasoft.io
// 📝 Github Repository: https://github.com/khulnasoft/velocity

package main

import (
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Create new Velocity instance
	app := velocity.New()

	// Create new GET route on path "/hello"
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start server on http://localhost:3000
	log.Fatal(app.Listen(":3000"))
}
