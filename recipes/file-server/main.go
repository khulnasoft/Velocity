// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Static file server
	app.Static("/", "./files")
	// => http://localhost:3000/hello.txt
	// => http://localhost:3000/gopher.gif

	// Start server
	log.Fatal(app.Listen(":3000"))
}
