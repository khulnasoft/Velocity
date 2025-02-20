package routes

import (
	"github.com/zeimedee/go-postgres/database"
	"github.com/zeimedee/go-postgres/models"
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

	database.DB.Db.Create(&book)

	return c.Status(200).JSON(book)
}

// AllBooks
func AllBooks(c *velocity.Ctx) error {
	books := []models.Book{}
	database.DB.Db.Find(&books)

	return c.Status(200).JSON(books)
}

// Book
func Book(c *velocity.Ctx) error {
	book := []models.Book{}
	title := new(models.Book)
	if err := c.BodyParser(title); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Db.Where("title = ?", title.Title).Find(&book)
	return c.Status(200).JSON(book)
}

// Update
func Update(c *velocity.Ctx) error {
	book := []models.Book{}
	title := new(models.Book)
	if err := c.BodyParser(title); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DB.Db.Model(&book).Where("title = ?", title.Title).Update("author", title.Author)

	return c.Status(200).JSON("updated")
}

// Delete
func Delete(c *velocity.Ctx) error {
	book := []models.Book{}
	title := new(models.Book)
	if err := c.BodyParser(title); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Db.Where("title = ?", title.Title).Delete(&book)

	return c.Status(200).JSON("deleted")
}
