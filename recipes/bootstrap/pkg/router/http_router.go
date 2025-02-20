package router

import (
	"github.com/kooroshh/velocity-boostrap/app/controllers"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/cors"
	"go.khulnasoft.com/velocity/middleware/csrf"
)

type HttpRouter struct{}

func (h HttpRouter) InstallRouter(app *velocity.App) {
	group := app.Group("", cors.New(), csrf.New())
	group.Get("/", controllers.RenderHello)
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{}
}
