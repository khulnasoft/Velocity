// 🚀 Velocity is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://docs.khulnasoft.io
// 📝 Github Repository: https://github.com/khulnasoft/velocity

package main

import (
	"fmt"
	"log"

	"github.com/khulnasoft/contrib/websocket"
	"go.khulnasoft.com/velocity"
)

func main() {
	app := velocity.New()

	// Optional middleware
	app.Use("/ws", func(c *velocity.Ctx) error {
		if c.Get("host") == "localhost:3000" {
			c.Locals("Host", "Localhost:3000")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		fmt.Println(c.Locals("Host")) // "Localhost:3000"
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	// ws://localhost:3000/ws
	log.Fatal(app.Listen(":3000"))
}
