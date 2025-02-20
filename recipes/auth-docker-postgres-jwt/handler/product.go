package handler

import (
	"app/database"
	"app/model"

	"go.khulnasoft.com/velocity"
)

// GetAllProducts query all products
func GetAllProducts(c *velocity.Ctx) error {
	db := database.DB
	var products []model.Product
	db.Find(&products)
	return c.JSON(velocity.Map{"status": "success", "message": "All products", "data": products})
}

// GetProduct query product
func GetProduct(c *velocity.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product
	db.Find(&product, id)
	if product.Title == "" {
		return c.Status(404).JSON(velocity.Map{"status": "error", "message": "No product found with ID", "data": nil})
	}
	return c.JSON(velocity.Map{"status": "success", "message": "Product found", "data": product})
}

// CreateProduct new product
func CreateProduct(c *velocity.Ctx) error {
	db := database.DB
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(500).JSON(velocity.Map{"status": "error", "message": "Couldn't create product", "data": err})
	}
	db.Create(&product)
	return c.JSON(velocity.Map{"status": "success", "message": "Created product", "data": product})
}

// DeleteProduct delete product
func DeleteProduct(c *velocity.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var product model.Product
	db.First(&product, id)
	if product.Title == "" {
		return c.Status(404).JSON(velocity.Map{"status": "error", "message": "No product found with ID", "data": nil})
	}
	db.Delete(&product)
	return c.JSON(velocity.Map{"status": "success", "message": "Product successfully deleted", "data": nil})
}
