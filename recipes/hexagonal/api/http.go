package api

import (
	"catalog/domain"
	"go.khulnasoft.com/velocity"
)

type handler struct {
	productService domain.Service
}

func (h *handler) Get(ctx *velocity.Ctx) error {
	code := ctx.Query("code")
	p, err := h.productService.Find(code)
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}
	return ctx.JSON(&p)
}

func (h *handler) Post(ctx *velocity.Ctx) error {
	p := &domain.Product{}
	if err := ctx.BodyParser(&p); err != nil {
		return ctx.Status(500).JSON(nil)
	}
	err := h.productService.Store(p)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	return ctx.JSON(&p)
}

func (h *handler) Put(ctx *velocity.Ctx) error {
	p := &domain.Product{}
	if err := ctx.BodyParser(&p); err != nil {
		return ctx.Status(500).JSON(nil)
	}
	err := h.productService.Update(p)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	return ctx.JSON(&p)
}

func (h *handler) Delete(ctx *velocity.Ctx) error {
	code := ctx.Params("code")
	err := h.productService.Delete(code)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	return ctx.Status(201).JSON(velocity.Map{"message": "ok"})
}

func (h *handler) GetAll(ctx *velocity.Ctx) error {
	p, err := h.productService.FindAll()
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}
	return ctx.JSON(&p)
}

// NewHandler  New handler instantiates a http handler for our product service
func NewHandler(productService domain.Service) *handler {
	return &handler{productService: productService}
}
