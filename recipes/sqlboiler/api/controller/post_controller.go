package controller

import (
	"context"
	"strconv"

	"velocity-sqlboiler/database"
	"velocity-sqlboiler/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.khulnasoft.com/velocity"
)

func GetPosts(c *velocity.Ctx) error {
	posts, err := models.Posts().All(context.Background(), database.DB)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(posts)
}

func GetPost(c *velocity.Ctx) error {
	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	post, err := models.FindPost(context.Background(), database.DB, postId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON(post)
}

func NewPost(c *velocity.Ctx) error {
	post := models.Post{}
	if err := c.BodyParser(&post); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if err := post.Insert(context.Background(), database.DB, boil.Infer()); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(post)
}

func DeletePost(c *velocity.Ctx) error {
	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	post, err := models.FindPost(context.Background(), database.DB, postId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}
	if _, err := post.Delete(context.Background(), database.DB); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.SendStatus(200)
}

func UpdatePost(c *velocity.Ctx) error {
	id := c.Params("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	newPost := models.Post{}
	if err := c.BodyParser(&newPost); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	post, err := models.FindPost(context.Background(), database.DB, postId)
	if err != nil {
		return c.Status(404).JSON(err.Error())
	}

	post.Title = newPost.Title
	post.Content = newPost.Content
	if _, err := post.Update(context.Background(), database.DB, boil.Infer()); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(post)
}
