package main

import (
	"log"

	"geoip-maxmind/handlers"

	"go.khulnasoft.com/velocity"
)

func main() {
	app := fiber.New()

	app.Get("/geo/:ip?", handlers.GeoIP)

	// Listen on port :3000
	log.Fatal(app.Listen(":3000"))
}
