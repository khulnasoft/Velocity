package main

import (
	"log"

	"fiber-sqlboiler/api/route"
	"fiber-sqlboiler/database"

	"go.khulnasoft.com/velocity"
)

func init() {
	database.ConnectDB()
}

func main() {
	app := fiber.New()
	route.SetupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
