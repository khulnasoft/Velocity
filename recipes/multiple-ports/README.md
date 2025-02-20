---
title: Multiple Ports
keywords: [multiple ports, server, port]
description: Running an application on multiple ports.
---

# Multiple Ports Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/multiple-ports) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/multiple-ports)

This project demonstrates how to run a Go application using the Velocity framework on multiple ports.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/multiple-ports
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

Here is an example of how to run a Velocity application on multiple ports:

```go
package main

import (
    "log"
    "sync"

    "go.khulnasoft.com/velocity"
)

func main() {
    app := velocity.New()

    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    ports := []string{":3000", ":3001"}

    var wg sync.WaitGroup
    for _, port := range ports {
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            if err := app.Listen(p); err != nil {
                log.Printf("Error starting server on port %s: %v", p, err)
            }
        }(port)
    }

    wg.Wait()
}
```

In this example:
- The application listens on multiple ports (`:3000` and `:3001`).
- A `sync.WaitGroup` is used to wait for all goroutines to finish.

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
