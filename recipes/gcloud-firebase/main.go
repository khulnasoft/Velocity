package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"go.khulnasoft.com/velocity"
)

var (
	app       *velocity.App
	fbApp     *firebase.App
	projectID string
)

// Hero db heroes struct
type Hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	projectID = os.Getenv("GCP_PROJECT")
	if projectID == "" {
		// App Engine uses another name
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}

	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://" + projectID + ".firebaseio.com",
	}

	fbApp, err = firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("functions.init: NewApp %v\n", err)
	}

	db, err := fbApp.Database(ctx)
	if err != nil {
		log.Fatalf("functions.init: Database init : %v\n", err)
	}

	heroesRef := db.NewRef("heroes")

	app = velocity.New()

	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("Health check ✅")
	})

	group := app.Group("api")

	group.Get("/hello", func(c *velocity.Ctx) error {
		return c.SendString("Hello World 🚀")
	})

	group.Get("/ola", func(c *velocity.Ctx) error {
		return c.SendString("Olá Mundo 🚀")
	})

	group.Get("/heroes", func(c *velocity.Ctx) error {
		ctx := context.Background()
		var heroes map[string]Hero
		if err := heroesRef.Get(ctx, &heroes); err != nil {
			return c.JSON(map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(map[string]map[string]Hero{
			"heroes": heroes,
		})
	})

	group.Get("/heroes/:id", func(c *velocity.Ctx) error {
		id := c.Params("id")
		var hero Hero
		if err := heroesRef.Child(id).Get(ctx, &hero); err != nil {
			return c.JSON(map[string]interface{}{
				"message": err.Error(),
			})
		}
		if hero.ID == "" {
			return c.Status(velocity.StatusNotFound).JSON(map[string]interface{}{
				"message": "Not Found",
			})
		}
		return c.JSON(hero)
	})
}

// Start start Velocity app with normal interface
func Start(addr string) error {
	return app.Listen(addr)
}

// HeroesAPI Exported http.HandlerFunc to be deployed to as a Cloud Function
func HeroesAPI(w http.ResponseWriter, r *http.Request) {
	err := CloudFunctionRouteToVelocity(app, w, r)
	if err != nil {
		fmt.Fprintf(w, "err : %v", err)
		return
	}
}
