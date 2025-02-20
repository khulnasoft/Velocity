---
id: error-handling
title: 🐛 Error Handling
description: >-
  Velocity supports centralized error handling by returning an error to the handler
  which allows you to log errors to external services or send a customized HTTP
  response to the client.
sidebar_position: 4
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Catching Errors

It’s essential to ensure that Velocity catches all errors that occur while running route handlers and middleware. You must return them to the handler function, where Velocity will catch and process them.

<Tabs>
<TabItem value="example" label="Example">

```go
app.Get("/", func(c velocity.Ctx) error {
    // Pass error to Velocity
    return c.SendFile("file-does-not-exist")
})
```

</TabItem>
</Tabs>

Velocity does not handle [panics](https://go.dev/blog/defer-panic-and-recover) by default. To recover from a panic thrown by any handler in the stack, you need to include the `Recover` middleware below:

```go title="Example"
package main

import (
    "log"

    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/middleware/recover"
)

func main() {
    app := velocity.New()

    app.Use(recover.New())

    app.Get("/", func(c velocity.Ctx) error {
        panic("This panic is caught by velocity")
    })

    log.Fatal(app.Listen(":3000"))
}
```

You could use Velocity's custom error struct to pass an additional `status code` using `velocity.NewError()`. It's optional to pass a message; if this is left empty, it will default to the status code message \(`404` equals `Not Found`\).

```go title="Example"
app.Get("/", func(c velocity.Ctx) error {
    // 503 Service Unavailable
    return velocity.ErrServiceUnavailable

    // 503 On vacation!
    return velocity.NewError(velocity.StatusServiceUnavailable, "On vacation!")
})
```

## Default Error Handler

Velocity provides an error handler by default. For a standard error, the response is sent as **500 Internal Server Error**. If the error is of type [velocity.Error](https://godoc.org/github.com/khulnasoft/velocity#Error), the response is sent with the provided status code and message.

```go title="Example"
// Default error handler
var DefaultErrorHandler = func(c velocity.Ctx, err error) error {
    // Status code defaults to 500
    code := velocity.StatusInternalServerError

    // Retrieve the custom status code if it's a *velocity.Error
    var e *velocity.Error
    if errors.As(err, &e) {
        code = e.Code
    }

    // Set Content-Type: text/plain; charset=utf-8
    c.Set(velocity.HeaderContentType, velocity.MIMETextPlainCharsetUTF8)

    // Return status code with error message
    return c.Status(code).SendString(err.Error())
}
```

## Custom Error Handler

A custom error handler can be set using a [Config](../api/velocity.md#errorhandler) when initializing a [Velocity instance](../api/velocity.md#new).

In most cases, the default error handler should be sufficient. However, a custom error handler can come in handy if you want to capture different types of errors and take action accordingly e.g., send a notification email or log an error to the centralized system. You can also send customized responses to the client e.g., error page or just a JSON response.

The following example shows how to display error pages for different types of errors.

```go title="Example"
// Create a new velocity instance with custom config
app := velocity.New(velocity.Config{
    // Override default error handler
    ErrorHandler: func(ctx velocity.Ctx, err error) error {
        // Status code defaults to 500
        code := velocity.StatusInternalServerError

        // Retrieve the custom status code if it's a *velocity.Error
        var e *velocity.Error
        if errors.As(err, &e) {
            code = e.Code
        }

        // Send custom error page
        err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
        if err != nil {
            // In case the SendFile fails
            return ctx.Status(velocity.StatusInternalServerError).SendString("Internal Server Error")
        }

        // Return from handler
        return nil
    },
})

// ...
```

> Special thanks to the [Echo](https://echo.labstack.com/) & [Express](https://expressjs.com/) framework for inspiration regarding error handling.
