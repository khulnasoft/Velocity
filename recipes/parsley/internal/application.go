package internal

import (
	"context"

	"parsley-app/internal/route_handlers"

	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"go.khulnasoft.com/velocity"
)

type parsleyApplication struct {
	app *velocity.App
}

var _ bootstrap.Application = &parsleyApplication{}

// NewApp Creates the main application service instance. This constructor function gets invoked by Parsley; add parameters for all required services.
func NewApp(app *velocity.App, routeHandlers []route_handlers.RouteHandler) bootstrap.Application {
	// Register RouteHandler services with the resolved Velocity instance.
	for _, routeHandler := range routeHandlers {
		routeHandler.Register(app)
	}

	return &parsleyApplication{
		app: app,
	}
}

// Run The entrypoint for the Parsley application.
func (a *parsleyApplication) Run(_ context.Context) error {
	return a.app.Listen(":5502")
}
