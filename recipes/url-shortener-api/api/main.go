// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amalshaji/stoyle/routes"
	"github.com/joho/godotenv"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

// setup two routes, one for shortening the url
// the other for resolving the url
// for example if the short is `4fg`, the user
// must navigate to `localhost:3000/4fg` to redirect to
// original URL. The domain can be changes in .env file
func setupRoutes(app *velocity.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	app := velocity.New()

	// app.Use(csrf.New())
	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
