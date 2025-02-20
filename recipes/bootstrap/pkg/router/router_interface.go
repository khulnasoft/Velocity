package router

import "go.khulnasoft.com/velocity"

type Router interface {
	InstallRouter(app *fiber.App)
}
