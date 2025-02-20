package controller

import (
	"context"
	"strconv"
	"time"

	"velocity-sqlc/database"
	"velocity-sqlc/database/sqlc"

	"go.khulnasoft.com/velocity"
)

func GetAuthors(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	authors, err := sqlc.New(database.DB).GetAuthors(ctx)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(authors)
}

func GetAuthor(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	author, err := sqlc.New(database.DB).GetAuthor(ctx, int32(authorId))
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(author)
}

func NewAuthor(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	var author sqlc.NewAuthorParams
	if err := c.BodyParser(&author); err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	newAuthor, err := sqlc.New(database.DB).NewAuthor(ctx, author)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(newAuthor)
}

func DeleteAuthor(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	err = sqlc.New(database.DB).DeleteAuthor(ctx, int32(authorId))
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).SendString("Author successfully deleted")
}

func UpdateAuthor(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	var author sqlc.UpdateAuthorParams
	author.ID = int32(authorId)
	if err := c.BodyParser(&author); err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	updatedAuthor, err := sqlc.New(database.DB).UpdateAuthor(ctx, author)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(updatedAuthor)
}
