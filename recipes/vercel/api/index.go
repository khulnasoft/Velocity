package handler

import (
	"net/http"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/adaptor"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	// This is needed to set the proper request path in `*velocity.Ctx`
	r.RequestURI = r.URL.String()

	handler().ServeHTTP(w, r)
}

// building the velocity application
func handler() http.HandlerFunc {
	app := velocity.New()

	app.Get("/v1", func(ctx *velocity.Ctx) error {
		return ctx.JSON(velocity.Map{
			"version": "v1",
		})
	})

	app.Get("/v2", func(ctx *velocity.Ctx) error {
		return ctx.JSON(velocity.Map{
			"version": "v2",
		})
	})

	app.Get("/", func(ctx *velocity.Ctx) error {
		return ctx.JSON(velocity.Map{
			"uri":  ctx.Request().URI().String(),
			"path": ctx.Path(),
		})
	})

	return adaptor.VelocityApp(app)
}
