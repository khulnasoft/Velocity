package router

import (
	"go.khulnasoft.com/velocity"
)

func InstallRouter(app *velocity.App) {
	setup(app, NewApiRouter(), NewHttpRouter())
}

func setup(app *velocity.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
