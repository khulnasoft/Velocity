package main

import (
	"fmt"
	"log"

	"github.com/kooroshh/velocity-boostrap/bootstrap"
	"github.com/kooroshh/velocity-boostrap/pkg/env"
)

func main() {
	app := bootstrap.NewApplication()
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT", "4000"))))
}
