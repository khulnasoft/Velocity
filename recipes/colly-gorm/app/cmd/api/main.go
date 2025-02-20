package main

import (
	"log"

	"velocity-colly-gorm/internals/consts"
	"velocity-colly-gorm/internals/services/database"
	"velocity-colly-gorm/internals/services/scrapers"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
	"go.khulnasoft.com/velocity/middleware/logger"
)

func main() {
	config, err := consts.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables!\n", err.Error())
	}
	database.ConnectDb(&config)

	app := velocity.New()
	micro := velocity.New()
	scrape := velocity.New()

	app.Mount("/api", micro)
	app.Mount("/scrape", scrape)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET",
		AllowCredentials: true,
	}))

	micro.Get("/healthchecker", func(c *velocity.Ctx) error {
		return c.Status(200).JSON(velocity.Map{
			"status":  "success",
			"message": "Welcome to Golang, Velocity, and Colly",
		})
	})

	scrape.Get("quotes", func(c *velocity.Ctx) error {
		go scrapers.Quotes()
		return c.Status(200).JSON(velocity.Map{
			"status":  "success",
			"message": "Start scraping quotes.toscrape.com ...",
		})
	})

	scrape.Get("coursera", func(c *velocity.Ctx) error {
		go scrapers.CourseraCourses()
		return c.Status(200).JSON(velocity.Map{
			"status":  "success",
			"message": "Start scraping courses details from coursera.org...",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
