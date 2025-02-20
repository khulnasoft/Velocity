package main

import (
	"log"

	"velocity-sqlc/api/route"
	"velocity-sqlc/database"

	"go.khulnasoft.com/velocity"
)

func init() {
	database.ConnectDB()
}

func main() {
	app := velocity.New()
	route.SetupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
