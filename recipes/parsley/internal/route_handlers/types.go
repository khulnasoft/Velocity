package route_handlers

import "github.com/khulnasoft/fiber/v2"

// RouteHandler Must be implemented by route handler services.
type RouteHandler interface {
	Register(app *fiber.App)
}
