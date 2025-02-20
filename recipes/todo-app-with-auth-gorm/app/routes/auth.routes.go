package routes

import (
	"numtostr/gotodo/app/services"

	"go.khulnasoft.com/velocity"
)

// AuthRoutes containes all the auth routes
func AuthRoutes(app velocity.Router) {
	r := app.Group("/auth")

	r.Post("/signup", services.Signup)
	r.Post("/login", services.Login)
}
