package main

// thanks to https://github.com/Learn-by-doing/csrf-examples
import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/khulnasoft/template/html/v2"
	"go.khulnasoft.com/velocity"
	"main/routes"
)

//go:embed views/*
var viewsfs embed.FS

func main() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	go func() {
		// ### EVIL SERVER ###
		// Velocity instance
		app := velocity.New(velocity.Config{Views: engine})
		app.Get("/", func(c *velocity.Ctx) error {
			// Render index - start with views directory
			return c.Render("views/layouts/main", velocity.Map{
				"Title": "Hello, World!",
			})
		})
		routes.RegisterEvilRoutes(app)
		fmt.Println("\"Evil\" server started and listening at localhost:3001")
		// Start server
		log.Fatal(app.Listen(":3001"))
	}()

	// ### NORMAL SERVER ###
	// Velocity instance
	app := velocity.New(velocity.Config{Views: engine})
	app.Get("/", func(c *velocity.Ctx) error {
		// Render index - start with views directory
		return c.Render("views/layouts/main", velocity.Map{
			"Title": "Hello, World!",
		})
	})
	routes.RegisterRoutes(app)
	fmt.Printf("Server started and listening at localhost:3000 - csrfActive: %v\n", len(os.Args) > 1 && os.Args[1] == "withoutCsrf")
	// Start server
	log.Fatal(app.Listen(":3000"))
}
