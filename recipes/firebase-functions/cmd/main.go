package main

import (
	"log"
	"os"

	app "example.com/KhulnasoftFirebaseBoilerplate"
)

func main() {
	port := "3001"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := app.Start(port); err != nil {
		log.Fatalf("app.Start: %v\n", err)
	}
}
