package controllers

import "go.khulnasoft.com/velocity"

func RenderHello(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"FiberTitle": "Hello From Fiber Html Engine",
	})
}
