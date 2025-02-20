---
title: 404 Handler
keywords: [404, not found, handler, errorhandler, custom]
description: Custom 404 error page handling.
---

# Custom 404 Not Found Handler Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/404-handler) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/404-handler)

This example demonstrates how to implement a custom 404 Not Found handler using the [Velocity](https://khulnasoft.io) web framework in Go. The purpose of this example is to show how to handle requests to undefined routes gracefully by returning a 404 status code.

## Description

In web applications, it's common to encounter requests to routes that do not exist. Handling these requests properly is important to provide a good user experience and to inform the user that the requested resource is not available. This example sets up a simple Velocity application with a custom 404 handler to manage such cases.

## Requirements

- [Go](https://golang.org/dl/) 1.18 or higher
- [Git](https://git-scm.com/downloads)

## Running the Example

To run the example, use the following command:
```bash
go run main.go
```

The server will start and listen on `localhost:3000`.

## Example Routes

- **GET /hello**: Returns a simple greeting message.
- **Undefined Routes**: Any request to a route not defined will trigger the custom 404 handler.

## Custom 404 Handler

The custom 404 handler is defined to catch all undefined routes and return a 404 status code with a "Not Found" message.

## Code Overview

### `main.go`

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
)

func main() {
    // Velocity instance
    app := velocity.New()

    // Routes
    app.Get("/hello", hello)

    // 404 Handler
    app.Use(func(c *velocity.Ctx) error {
        return c.SendStatus(404) // => 404 "Not Found"
    })

    // Start server
    log.Fatal(app.Listen(":3000"))
}

// Handler
func hello(c *velocity.Ctx) error {
    return c.SendString("I made a ☕ for you!")
}
```

## Conclusion

This example provides a basic setup for handling 404 Not Found errors in a Velocity application. It can be extended and customized further to fit the needs of more complex applications.

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [GitHub Repository](https://github.com/khulnasoft/velocity)
