package main

import (
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	app := velocity.New()

	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
