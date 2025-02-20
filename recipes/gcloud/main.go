package app

import (
	"fmt"
	"net/http"
	"strings"

	"go.khulnasoft.com/velocity"
)

var app *velocity.App

func init() {
	app = velocity.New()

	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Health check ✅")
	})

	group := app.Group("api")

	group.Get("/hello", func(c *velocity.Ctx) error {
		return c.SendString("Hello World 🚀")
	})

	group.Get("/ola", func(c *velocity.Ctx) error {
		return c.SendString("Olá Mundo 🚀")
	})
}

// Start start Velocity app with normal interface
func Start(addr string) error {
	if -1 == strings.IndexByte(addr, ':') {
		addr = ":" + addr
	}

	return app.Listen(addr)
}

// MyCloudFunction Exported http.HandlerFunc to be deployed to as a Cloud Function
func MyCloudFunction(w http.ResponseWriter, r *http.Request) {
	err := CloudFunctionRouteToVelocity(app, w, r)
	if err != nil {
		fmt.Fprintf(w, "err : %v", err)
		return
	}
}
