package main

import (
	"log"

	"api-velocity-gorm/database"
	"api-velocity-gorm/router"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
)

func main() {
	app := velocity.New()
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
