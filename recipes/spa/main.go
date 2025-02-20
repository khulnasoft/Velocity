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

	// serve Single Page application on "/web"
	// assume static file at dist folder
	app.Static("/web", "dist")

	app.Get("/web/*", func(ctx *velocity.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})

	// Start server on http://localhost:3000
	log.Fatal(app.Listen(":3000"))
}
