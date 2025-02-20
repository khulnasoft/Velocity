// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/fiber
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"go.khulnasoft.com/velocity"
	"github.com/joho/godotenv"
	khulnasoftfirebaseauth "github.com/sacsand/khulnasoft-firebaseauth"
	"google.golang.org/api/option"
)

func init() {
	// Loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Fiber instance
	app := fiber.New()

	// Get google service account credentials
	serviceAccount, fileExi := os.LookupEnv("GOOGLE_SERVICE_ACCOUNT")

	if !fileExi {
		log.Fatalf("Please provide valid firebbase auth credential json!")
	}

	// Initialize the firebase app.
	opt := option.WithCredentialsFile(serviceAccount)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// Unauthenticated routes
	app.Get("/salaanthe", salanthe)

	// Initialize the middleware with config. See https://github.com/sacsand/khulnasoft-firebaseauth for more configuration options.
	app.Use(khulnasoftfirebaseauth.New(khulnasoftfirebaseauth.Config{
		// Firebase Authentication App Object
		// Mandatory
		FirebaseApp: fireApp,
		// Ignore urls array.
		// Optional. These url will ignore by middleware
		IgnoreUrls: []string{"GET::/salut", "POST::/ciao"},
	}))

	// Authenticaed Routes.
	app.Get("/hello", hello)
	app.Get("/salut", salut) // Ignore the auth by IgnoreUrls config
	app.Post("/ciao", ciao)  // Ignore the auth by IgnoreUrls config
	app.Get("/ayubowan", ayubowan)

	// Start server.
	log.Fatal(app.Listen(":3001"))
}

/**
*
* Controllers
*
 */

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func salut(c *fiber.Ctx) error {
	return c.SendString("Salut, World 👋!")
}

func ciao(c *fiber.Ctx) error {
	return c.SendString("Ciao, World 👋! ")
}

func ayubowan(c *fiber.Ctx) error {
	// Get authenticated user from context.
	claims := c.Locals("user")
	fmt.Println(claims)
	return c.SendString("Ayubowan👋! ")
}

func salanthe(c *fiber.Ctx) error {
	return c.SendString("Salanthe👋! ")
}
