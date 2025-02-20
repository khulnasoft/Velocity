package main

import (
	"log"

	"app/database"
	"app/router"

	"go.khulnasoft.com/velocity"
	// "go.khulnasoft.com/velocity/middleware/cors"
)

func main() {
	app := velocity.New(velocity.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Velocity",
		AppName:       "App Name",
	})
	// app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
