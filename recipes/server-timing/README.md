---
title: Server Timing
keywords: [server timing]
description: Adding Server Timing headers to an application.
---

# Server Timing

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/server-timing) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/server-timing)

This project demonstrates how to implement Server-Timing headers in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/server-timing
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

## Example

Here is an example of how to set up Server-Timing headers in a Velocity application:

```go
package main

import (
    "go.khulnasoft.com/velocity"
    "time"
)

func main() {
    app := velocity.New()

    app.Use(func(c *velocity.Ctx) error {
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        c.Append("Server-Timing", "app;dur="+duration.String())
        return err
    })

    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Server-Timing Header Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server-Timing)
