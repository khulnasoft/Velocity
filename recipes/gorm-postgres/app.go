package main

import (
	"log"

	"github.com/zeimedee/go-postgres/database"
	"github.com/zeimedee/go-postgres/routes"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
)

func setUpRoutes(app *velocity.App) {
	app.Get("/hello", routes.Hello)
	app.Get("/allbooks", routes.AllBooks)
	app.Post("/addbook", routes.AddBook)
	app.Post("/book", routes.Book)
	app.Put("/update", routes.Update)
	app.Delete("/delete", routes.Delete)
}

func main() {
	database.ConnectDb()
	app := velocity.New()

	setUpRoutes(app)

	app.Use(cors.New())

	app.Use(func(c *velocity.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":3000"))
}
