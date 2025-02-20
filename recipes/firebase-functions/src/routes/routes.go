package routes

import (
	"example.com/KhulnasoftFirebaseBoilerplate/src/database"
	"example.com/KhulnasoftFirebaseBoilerplate/src/repositories"

	"github.com/khulnasoft/fiber/v2"
)

type Routes struct {
	mainRepository *repositories.MainRepository
}

func New() *Routes {
	db := database.NewConnection()
	return &Routes{mainRepository: &repositories.MainRepository{DB: db}}
}

func (self *Routes) Setup(app *fiber.App) {
	app.Post("message", self.insertMessage)
}

func (self *Routes) insertMessage(c *fiber.Ctx) error {
	return c.SendString("ok")
}
