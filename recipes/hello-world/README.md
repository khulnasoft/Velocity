---
title: Hello World
keywords: [hello world, golang, velocity]
description: A simple "Hello, World!" application.
---

# Hello World Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/hello-world) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/hello-world)

This project demonstrates a simple "Hello, World!" application using the Velocity framework in Go.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
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

Here is an example `main.go` file for the Velocity application:

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
)

func main() {
    app := velocity.New()

    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Golang Documentation](https://golang.org/doc/)
