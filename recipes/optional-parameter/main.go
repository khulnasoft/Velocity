// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"log"
	"strconv"

	"go.khulnasoft.com/velocity"
)

func main() {
	// user list
	users := [...]string{"Alice", "Bob", "Charlie", "David"}

	// Velocity instance
	app := velocity.New()

	// Route to profile
	app.Get("/:id?", func(c *velocity.Ctx) error {
		id, err := strconv.Atoi(c.Params("id")) // transform id to array index
		if err != nil || id < 0 || id >= len(users) {
			return c.SendStatus(404) // invalid parameter returns 404
		}
		return c.SendString("Hello, " + users[id] + "!") // custom hello message to user with the id
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
