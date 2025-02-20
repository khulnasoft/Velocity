package api

import (
	"go.khulnasoft.com/velocity"
)

// ProductHandler  an interface with operations to be implemented by a specific handler, ie http, gRCP
type ProductHandler interface {
	Get(ctx *velocity.Ctx)
	Post(ctx *velocity.Ctx)
	Put(ctx *velocity.Ctx)
	Delete(ctx *velocity.Ctx)
	GetAll(ctx *velocity.Ctx)
}
