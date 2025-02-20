// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"fmt"
	"log"
	"time"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Custom Timer middleware
	app.Use(Timer())

	// Routes
	app.Get("/", func(c *velocity.Ctx) error {
		time.Sleep(2 * time.Second) // Sleep 2 seconds
		return c.SendString("That took a while 😞")
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Timer will measure how long it takes before a response is returned
func Timer() velocity.Handler {
	return func(c *velocity.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		// Do something with response
		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))
		// return stack error if exist
		return err
	}
}
