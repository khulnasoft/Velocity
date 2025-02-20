package api

import (
	"go.khulnasoft.com/velocity"
)

// ProductHandler  an interface with operations to be implemented by a specific handler, ie http, gRCP
type ProductHandler interface {
	Get(ctx *fiber.Ctx)
	Post(ctx *fiber.Ctx)
	Put(ctx *fiber.Ctx)
	Delete(ctx *fiber.Ctx)
	GetAll(ctx *fiber.Ctx)
}
