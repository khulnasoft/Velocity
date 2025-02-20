package controller

import (
	"context"
	"strconv"
	"time"

	"velocity-sqlc/database"
	"velocity-sqlc/database/sqlc"

	"go.khulnasoft.com/velocity"
)

func GetPosts(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	posts, err := sqlc.New(database.DB).GetPosts(ctx)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(posts)
}

func GetPost(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	post, err := sqlc.New(database.DB).GetPost(ctx, int32(postId))
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(post)
}

func NewPost(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	var post sqlc.NewPostParams
	if err := c.BodyParser(&post); err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	newPost, err := sqlc.New(database.DB).NewPost(ctx, post)
	if err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).JSON(newPost)
}

func DeletePost(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	if err := sqlc.New(database.DB).DeletePost(ctx, int32(postId)); err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).SendString("Post deleted")
}

func UpdatePost(c *velocity.Ctx) error {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	var post sqlc.UpdatePostParams
	if err := c.BodyParser(&post); err != nil {
		return c.Status(velocity.StatusBadRequest).SendString(err.Error())
	}

	post.ID = int32(postId)
	if _, err := sqlc.New(database.DB).UpdatePost(ctx, post); err != nil {
		return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(velocity.StatusOK).SendString("Post updated")
}
