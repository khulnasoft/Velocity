package main

import (
	"context"
	"time"

	"github.com/amalshaji/velocity-netlify/adapter"
	"github.com/amalshaji/velocity-netlify/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.khulnasoft.com/velocity"
)

var velocityLambda *adapter.VelocityLambda

func init() {
	app := velocity.New()
	app.Static("/", "./public")
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendFile("index")
	})
	app.Get("/api/:ip", handler.CacheRequest(10*time.Minute), handler.GeoLocation)

	velocityLambda = adapter.New(app)
}

// Handler proxies our app requests to aws lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return velocityLambda.ProxyWithContext(ctx, req)
}

func main() {
	//
	lambda.Start(Handler)
}
