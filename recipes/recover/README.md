---
title: Recover Middleware
keywords: [recover, middleware]
description: Recover middleware for error handling.
---

# Recover Middleware Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/recover) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/recover)

This project demonstrates how to implement a recovery mechanism in a Go application using the Velocity framework's `Recover` middleware.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/recover
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

Here is an example of how to set up the `Recover` middleware in a Velocity application:

```go
package main

import (
    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/recover"
)

func main() {
    app := velocity.New()

    // Use the Recover middleware
    app.Use(recover.New())

    app.Get("/", func(c *velocity.Ctx) error {
        // This will cause a panic
        panic("something went wrong")
    })

    app.Listen(":3000")
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Velocity Recover Middleware Documentation](https://docs.khulnasoft.io/api/middleware/recover)
