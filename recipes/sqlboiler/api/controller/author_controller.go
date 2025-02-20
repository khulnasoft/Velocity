package controller

import (
	"context"
	"strconv"

	"velocity-sqlboiler/database"
	"velocity-sqlboiler/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.khulnasoft.com/velocity"
)

func GetAuthors(c *velocity.Ctx) error {
	authors, err := models.Authors().All(context.Background(), database.DB)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(authors)
}

func GetAuthor(c *velocity.Ctx) error {
	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	author, err := models.FindAuthor(context.Background(), database.DB, authorId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON(author)
}

func NewAuthor(c *velocity.Ctx) error {
	author := models.Author{}
	if err := c.BodyParser(&author); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if err := author.Insert(context.Background(), database.DB, boil.Infer()); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(author)
}

func DeleteAuthor(c *velocity.Ctx) error {
	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	author, err := models.FindAuthor(context.Background(), database.DB, authorId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}
	if _, err := author.Delete(context.Background(), database.DB); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.SendStatus(200)
}

func UpdateAuthor(c *velocity.Ctx) error {
	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	newAuthor := models.Author{}
	if err := c.BodyParser(&newAuthor); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	author, err := models.FindAuthor(context.Background(), database.DB, authorId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}

	author.Name = newAuthor.Name
	if _, err := author.Update(context.Background(), database.DB, boil.Infer()); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(author)
}
