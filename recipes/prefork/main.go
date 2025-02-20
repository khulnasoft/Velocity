// 🚀 Velocity is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://docs.khulnasoft.io
// 📝 Github Repository: https://github.com/khulnasoft/velocity

package main

import (
	"fmt"
	"log"
	"os"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Print current process
	if velocity.IsChild() {
		fmt.Printf("[%d] Child\n", os.Getppid())
	} else {
		fmt.Printf("[%d] Master\n", os.Getppid())
	}

	// Velocity instance
	app := velocity.New(velocity.Config{
		Prefork: true,
	})

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(":3000"))

	// Run the following command to see all processes sharing port 3000:
	// sudo lsof -i -P -n | grep LISTEN
}

// Handler
func hello(c *velocity.Ctx) error {
	return c.SendString("Hello, World!")
}
