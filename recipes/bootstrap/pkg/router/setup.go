package router

import (
	"github.com/khulnasoft/fiber/v2"
)

func InstallRouter(app *fiber.App) {
	setup(app, NewApiRouter(), NewHttpRouter())
}

func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
