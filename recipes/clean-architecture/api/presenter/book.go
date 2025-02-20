package presenter

import (
	"clean-architecture/pkg/entities"
	"go.khulnasoft.com/velocity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book is the presenter object which will be passed in the response by Handler
type Book struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
}

// BookSuccessResponse is the singular SuccessResponse that will be passed in the response by
// Handler
func BookSuccessResponse(data *entities.Book) *velocity.Map {
	book := Book{
		ID:     data.ID,
		Title:  data.Title,
		Author: data.Author,
	}
	return &velocity.Map{
		"status": true,
		"data":   book,
		"error":  nil,
	}
}

// BooksSuccessResponse is the list SuccessResponse that will be passed in the response by Handler
func BooksSuccessResponse(data *[]Book) *velocity.Map {
	return &velocity.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

// BookErrorResponse is the ErrorResponse that will be passed in the response by Handler
func BookErrorResponse(err error) *velocity.Map {
	return &velocity.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
