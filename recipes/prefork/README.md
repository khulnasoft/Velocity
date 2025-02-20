---
title: Prefork
keywords: [prefork]
description: Running an application in prefork mode.
---

# Prefork Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/prefork) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/prefork)

This project demonstrates how to use the `Prefork` feature in a Go application using the Velocity framework. Preforking can improve performance by utilizing multiple CPU cores.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/prefork
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

Here is an example of how to set up the `Prefork` feature in a Velocity application:

```go
package main

import (
    "log"

    "go.khulnasoft.com/velocity"
)

func main() {
    // Velocity instance with Prefork enabled
    app := velocity.New(velocity.Config{
        Prefork: true,
    })

    // Routes
    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Velocity Prefork Documentation](https://docs.khulnasoft.io/api/velocity#prefork)
