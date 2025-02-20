package main

import (
	"velocity-oauth-google/router"

	"github.com/joho/godotenv"
	"go.khulnasoft.com/velocity"
)

func main() {
	godotenv.Load()
	app := velocity.New()
	router.Routes(app)
	app.Listen(":3300")
}
