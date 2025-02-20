package infrastructure

import (
	"fmt"
	"log"

	"docker-mariadb-clean-arch/internal/auth"
	"docker-mariadb-clean-arch/internal/city"
	"docker-mariadb-clean-arch/internal/misc"
	"docker-mariadb-clean-arch/internal/user"

	_ "github.com/go-sql-driver/mysql"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/compress"
	"go.khulnasoft.com/velocity/middleware/cors"
	"go.khulnasoft.com/velocity/middleware/etag"
	"go.khulnasoft.com/velocity/middleware/favicon"
	"go.khulnasoft.com/velocity/middleware/limiter"
	"go.khulnasoft.com/velocity/middleware/logger"
	"go.khulnasoft.com/velocity/middleware/recover"
	"go.khulnasoft.com/velocity/middleware/requestid"
)

// Run our Velocity webserver.
func Run() {
	// Try to connect to our database as the initial part.
	mariadb, err := ConnectToMariaDB()
	if err != nil {
		log.Fatal("Database connection error: $s", err)
	}

	// Creates a new Velocity instance.
	app := velocity.New(velocity.Config{
		AppName:      "Docker MariaDB Clean Arch",
		ServerHeader: "Velocity",
	})

	// Use global middlewares.
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(limiter.New(limiter.Config{
		Max: 100,
		LimitReached: func(c *velocity.Ctx) error {
			return c.Status(velocity.StatusTooManyRequests).JSON(&velocity.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())

	// Create repositories.
	cityRepository := city.NewCityRepository(mariadb)
	userRepository := user.NewUserRepository(mariadb)

	// Create all of our services.
	cityService := city.NewCityService(cityRepository)
	userService := user.NewUserService(userRepository)

	// Prepare our endpoints for the API.
	misc.NewMiscHandler(app.Group("/api/v1"))
	auth.NewAuthHandler(app.Group("/api/v1/auth"))
	city.NewCityHandler(app.Group("/api/v1/cities"), cityService)
	user.NewUserHandler(app.Group("/api/v1/users"), userService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *velocity.Ctx) error {
		errorMessage := fmt.Sprintf("Route '%s' does not exist in this API!", c.OriginalURL())

		return c.Status(velocity.StatusNotFound).JSON(&velocity.Map{
			"status":  "fail",
			"message": errorMessage,
		})
	})

	// Listen to port 8080.
	log.Fatal(app.Listen(":8080"))
}
