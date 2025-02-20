package main

import (
	"context"
	"log"

	"go.khulnasoft.com/velocity"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	velocityAdapter "github.com/awslabs/aws-lambda-go-api-proxy/velocity"
)

var velocityLambda *velocityAdapter.VelocityLambda

// init the Velocity Server
func init() {
	log.Printf("Velocity cold start")
	app := velocity.New()

	// Routes
	app.Get("/", func(c *velocity.Ctx) error {
		return c.JSON(velocity.Map{"message": "Hello World"})
	})

	velocityLambda = velocityAdapter.New(app)
}

// Handler will deal with Velocity working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return velocityLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
