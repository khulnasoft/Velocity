package controllers

import "go.khulnasoft.com/velocity"

func RenderHello(c *velocity.Ctx) error {
	return c.Render("index", velocity.Map{
		"VelocityTitle": "Hello From Velocity Html Engine",
	})
}
