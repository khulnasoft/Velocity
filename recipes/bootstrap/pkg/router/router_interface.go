package router

import "github.com/khulnasoft/fiber/v2"

type Router interface {
	InstallRouter(app *fiber.App)
}
