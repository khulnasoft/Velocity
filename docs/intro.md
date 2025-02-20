---
slug: /
id: welcome
title: 👋 Welcome
sidebar_position: 1
---
Welcome to the online API documentation for Velocity, complete with examples to help you start building web applications with Velocity right away!

**Velocity** is an [Express](https://github.com/expressjs/express)-inspired **web framework** built on top of [Fasthttp](https://github.com/valyala/fasthttp), the **fastest** HTTP engine for [Go](https://go.dev/doc/). It is designed to facilitate rapid development with **zero memory allocations** and a strong focus on **performance**.

These docs are for **Velocity v3**, which was released on **Month xx, 202x**.

### Installation

First, [download](https://go.dev/dl/) and install Go. Version `1.23` or higher is required.

Installation is done using the [`go get`](https://pkg.go.dev/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them) command:

```bash
go get go.khulnasoft.com/velocity/v3
```

### Zero Allocation

Velocity is optimized for **high performance**, meaning values returned from **velocity.Ctx** are **not** immutable by default and **will** be reused across requests. As a rule of thumb, you **must** only use context values within the handler and **must not** keep any references. Once you return from the handler, any values obtained from the context will be reused in future requests. Here is an example:

```go
func handler(c velocity.Ctx) error {
    // Variable is only valid within this handler
    result := c.Params("foo") 

    // ...
}
```

If you need to persist such values outside the handler, make copies of their **underlying buffer** using the [copy](https://pkg.go.dev/builtin/#copy) builtin. Here is an example for persisting a string:

```go
func handler(c velocity.Ctx) error {
    // Variable is only valid within this handler
    result := c.Params("foo")

    // Make a copy
    buffer := make([]byte, len(result))
    copy(buffer, result)
    resultCopy := string(buffer) 
    // Variable is now valid indefinitely

    // ...
}
```

We created a custom `CopyString` function that performs the above and is available under [khulnasoft/utils](https://github.com/khulnasoft/utils).

```go
app.Get("/:foo", func(c velocity.Ctx) error {
    // Variable is now immutable
    result := utils.CopyString(c.Params("foo")) 

    // ...
})
```

Alternatively, you can enable the `Immutable` setting. This makes all values returned from the context immutable, allowing you to persist them anywhere. Note that this comes at the cost of performance.

```go
app := velocity.New(velocity.Config{
    Immutable: true,
})
```

For more information, please refer to [#426](https://go.khulnasoft.com/velocity/issues/426), [#185](https://go.khulnasoft.com/velocity/issues/185), and [#3012](https://go.khulnasoft.com/velocity/issues/3012).

### Hello, World

Below is the most straightforward **Velocity** application you can create:

```go
package main

import "go.khulnasoft.com/velocity/v3"

func main() {
    app := velocity.New()

    app.Get("/", func(c velocity.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

```bash
go run server.go
```

Browse to `http://localhost:3000` and you should see `Hello, World!` displayed on the page.

### Basic Routing

Routing determines how an application responds to a client request to a particular endpoint, which is a URI (or path) and a specific HTTP request method (`GET`, `PUT`, `POST`, etc.).

Each route can have **multiple handler functions** that are executed when the route is matched.

Route definitions follow the structure below:

```go
// Function signature
app.Method(path string, ...func(velocity.Ctx) error)
```

- `app` is an instance of **Velocity**
- `Method` is an [HTTP request method](https://docs.khulnasoft.com/api/app#route-handlers): `GET`, `PUT`, `POST`, etc.
- `path` is a virtual path on the server
- `func(velocity.Ctx) error` is a callback function containing the [Context](https://docs.khulnasoft.com/api/ctx) executed when the route is matched

#### Simple Route

```go
// Respond with "Hello, World!" on root path "/"
app.Get("/", func(c velocity.Ctx) error {
    return c.SendString("Hello, World!")
})
```

#### Parameters

```go
// GET http://localhost:8080/hello%20world

app.Get("/:value", func(c velocity.Ctx) error {
    return c.SendString("value: " + c.Params("value"))
    // => Response: "value: hello world"
})
```

#### Optional Parameter

```go
// GET http://localhost:3000/john

app.Get("/:name?", func(c velocity.Ctx) error {
    if c.Params("name") != "" {
        return c.SendString("Hello " + c.Params("name"))
        // => Response: "Hello john"
    }
    return c.SendString("Where is john?")
    // => Response: "Where is john?"
})
```

#### Wildcards

```go
// GET http://localhost:3000/api/user/john

app.Get("/api/*", func(c velocity.Ctx) error {
    return c.SendString("API path: " + c.Params("*"))
    // => Response: "API path: user/john"
})
```

### Static Files

To serve static files such as **images**, **CSS**, and **JavaScript** files, use the [static middleware](./middleware/static.md).

Use the following code to serve files in a directory named `./public`:

```go
package main

import (
    "go.khulnasoft.com/velocity/v3"
    "go.khulnasoft.com/velocity/v3/middleware/static"
)

func main() {
    app := velocity.New()

    app.Use("/", static.New("./public"))

    app.Listen(":3000")
}
```

Now, you can access the files in the `./public` directory via your browser:

```bash
http://localhost:3000/hello.html
http://localhost:3000/js/jquery.js
http://localhost:3000/css/style.css
```
