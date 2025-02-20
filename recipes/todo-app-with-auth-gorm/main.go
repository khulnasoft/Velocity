package main

import (
	"fmt"

	"numtostr/gotodo/app/dal"
	"numtostr/gotodo/app/routes"
	"numtostr/gotodo/config"
	"numtostr/gotodo/config/database"
	"numtostr/gotodo/utils"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
)

func main() {
	database.Connect()
	database.Migrate(&dal.User{}, &dal.Todo{})

	app := velocity.New(velocity.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Use(logger.New())

	routes.AuthRoutes(app)
	routes.TodoRoutes(app)

	app.Listen(fmt.Sprintf(":%v", config.PORT))
}
