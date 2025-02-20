---
title: Hello World
keywords: [hello world, golang, fiber]
description: A simple "Hello, World!" application.
---

# Hello World Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://github.com/khulnasoft/recipes/tree/master/hello-world) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/hello-world)

This project demonstrates a simple "Hello, World!" application using the Fiber framework in Go.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Fiber](https://github.com/khulnasoft/fiber) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/khulnasoft/recipes.git
    cd recipes/hello-world
    ```

2. Install dependencies:
    ```sh
    go get
    ```

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. Access the application at `http://localhost:3000`.

## Example

Here is an example `main.go` file for the Fiber application:

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
)

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Fiber Documentation](https://docs.khulnasoft.io)
- [Golang Documentation](https://golang.org/doc/)
