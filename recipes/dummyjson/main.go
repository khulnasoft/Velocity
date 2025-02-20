package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"go.khulnasoft.com/velocity"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	app := velocity.New()

	app.Get("/", func(c *velocity.Ctx) error {
		resp, err := client.Get("https://dummyjson.com/products/1")
		if err != nil {
			return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).JSON(&velocity.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		if _, err := io.Copy(c.Response().BodyWriter(), resp.Body); err != nil {
			return c.Status(velocity.StatusInternalServerError).JSON(&velocity.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.SendStatus(velocity.StatusOK)
	})

	log.Fatal(app.Listen(":3000"))
}
