package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/khulnasoft/recipes/template-asset-bundling/handlers"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
	"go.khulnasoft.com/velocity/middleware/recover"
	"github.com/khulnasoft/template/html/v2"
)

func main() {
	// Create view engine
	engine := html.New("./views", ".html")

	// Disable this in production
	engine.Reload(true)

	engine.AddFunc("getCssAsset", func(name string) (res template.HTML) {
		filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
			}
			return nil
		})
		return
	})

	// Create fiber app
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Setup routes
	app.Get("/", handlers.Home)
	app.Get("/about", handlers.About)

	// Setup static files
	app.Static("/public", "./public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(":3000"))
}
