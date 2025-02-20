package bootstrap

import (
	"github.com/khulnasoft/template/html/v2"
	"github.com/kooroshh/velocity-boostrap/pkg/database"
	"github.com/kooroshh/velocity-boostrap/pkg/env"
	"github.com/kooroshh/velocity-boostrap/pkg/router"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
	"go.khulnasoft.com/velocity/middleware/monitor"
	"go.khulnasoft.com/velocity/middleware/recover"
)

func NewApplication() *velocity.App {
	env.SetupEnvFile()
	database.SetupDatabase()
	engine := html.New("./views", ".html")
	app := velocity.New(velocity.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	router.InstallRouter(app)

	return app
}
