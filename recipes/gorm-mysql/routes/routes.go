package routes

import (
	"strconv"

	"gorm-mysql/database"
	"gorm-mysql/models"

	"go.khulnasoft.com/velocity"
)

// Hello
func Hello(c *velocity.Ctx) error {
	return c.SendString("velocity")
}

// AddBook
func AddBook(c *velocity.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DBConn.Create(&book)

	return c.Status(200).JSON(book)
}

func GetBook(c *velocity.Ctx) error {
	books := []models.Book{}

	database.DBConn.First(&books, c.Params("id"))

	return c.Status(200).JSON(books)
}

// AllBooks
func AllBooks(c *velocity.Ctx) error {
	books := []models.Book{}

	database.DBConn.Find(&books)

	return c.Status(200).JSON(books)
}

// Update
func Update(c *velocity.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, _ := strconv.Atoi(c.Params("id"))

	database.DBConn.Model(&models.Book{}).Where("id = ?", id).Update("title", book.Title)

	return c.Status(200).JSON("updated")
}

// Delete
func Delete(c *velocity.Ctx) error {
	book := new(models.Book)

	id, _ := strconv.Atoi(c.Params("id"))

	database.DBConn.Where("id = ?", id).Delete(&book)

	return c.Status(200).JSON("deleted")
}
