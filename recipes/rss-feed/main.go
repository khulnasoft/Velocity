package main

import (
	"log"

	"github.com/khulnasoft/template/mustache/v2"
	"go.khulnasoft.com/velocity"
)

func main() {
	engineXML := mustache.New("./xmls", ".xml")
	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	app := velocity.New()

	app.Get("/rss", func(c *velocity.Ctx) error {
		// Set Content-Type to application/rss+xml
		c.Type("rss")

		// Set rendered template to body
		return engineXML.Render(c, "example", velocity.Map{
			"Lang":      "en",
			"Title":     "hello-rss",
			"Greetings": "Hello World",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
