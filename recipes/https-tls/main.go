// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"crypto/tls"
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Routes
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString(c.Protocol()) // => https
	})

	// Create tls certificate
	cer, err := tls.LoadX509KeyPair("certs/ssl.cert", "certs/ssl.key")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	// Create custom listener
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		panic(err)
	}

	// Start server with https/ssl enabled on http://localhost:443
	log.Fatal(app.Listener(ln))
}
