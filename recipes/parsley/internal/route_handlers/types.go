package route_handlers

import "go.khulnasoft.com/velocity"

// RouteHandler Must be implemented by route handler services.
type RouteHandler interface {
	Register(app *velocity.App)
}
