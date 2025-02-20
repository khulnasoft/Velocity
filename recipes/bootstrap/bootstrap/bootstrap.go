package bootstrap

import (
	"github.com/khulnasoft/fiber/v2"
	"github.com/khulnasoft/fiber/v2/middleware/logger"
	"github.com/khulnasoft/fiber/v2/middleware/monitor"
	"github.com/khulnasoft/fiber/v2/middleware/recover"
	"github.com/khulnasoft/template/html/v2"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"github.com/kooroshh/fiber-boostrap/pkg/router"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	database.SetupDatabase()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	router.InstallRouter(app)

	return app
}
