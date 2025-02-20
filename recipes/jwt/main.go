package main

import (
	"log"

	"api-velocity-gorm/database"
	"api-velocity-gorm/router"

	"github.com/joho/godotenv"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	app := velocity.New()
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
