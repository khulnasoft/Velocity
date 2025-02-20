package main

import (
	"fmt"
	"log"

	"velocity-gorm/book"
	"velocity-gorm/database"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *velocity.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("books.db"))
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	database.DBConn.AutoMigrate(&book.Book{})
	fmt.Println("Database Migrated")
}

func main() {
	app := velocity.New()
	app.Use(cors.New())

	initDatabase()

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
