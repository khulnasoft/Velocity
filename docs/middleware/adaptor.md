---
id: adaptor
---

# Adaptor

Converter for net/http handlers to/from Velocity request handlers, special thanks to [@arsmn](https://github.com/arsmn)!

## Signatures

| Name | Signature | Description
| :--- | :--- | :---
| HTTPHandler | `HTTPHandler(h http.Handler) velocity.Handler` | http.Handler -> velocity.Handler
| HTTPHandlerFunc | `HTTPHandlerFunc(h http.HandlerFunc) velocity.Handler` | http.HandlerFunc -> velocity.Handler
| HTTPMiddleware | `HTTPHandlerFunc(mw func(http.Handler) http.Handler) velocity.Handler` | func(http.Handler) http.Handler -> velocity.Handler
| VelocityHandler | `VelocityHandler(h velocity.Handler) http.Handler` | velocity.Handler -> http.Handler
| VelocityHandlerFunc | `VelocityHandlerFunc(h velocity.Handler) http.HandlerFunc` | velocity.Handler -> http.HandlerFunc
| VelocityApp | `VelocityApp(app *velocity.App) http.HandlerFunc` | Velocity app -> http.HandlerFunc
| ConvertRequest | `ConvertRequest(c velocity.Ctx, forServer bool) (*http.Request, error)` | velocity.Ctx -> http.Request
| CopyContextToVelocityContext | `CopyContextToVelocityContext(context any, requestContext *fasthttp.RequestCtx)` | context.Context -> fasthttp.RequestCtx

## Examples

### net/http to Velocity

```go
package main

import (
    "fmt"
    "net/http"

    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/adaptor"
)

func main() {
    // New velocity app
    app := velocity.New()

    // http.Handler -> velocity.Handler
    app.Get("/", adaptor.HTTPHandler(handler(greet)))

    // http.HandlerFunc -> velocity.Handler
    app.Get("/func", adaptor.HTTPHandlerFunc(greet))

    // Listen on port 3000
    app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
    return http.HandlerFunc(f)
}

func greet(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World!")
}
```

### net/http middleware to Velocity

```go
package main

import (
    "log"
    "net/http"

    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/adaptor"
)

func main() {
    // New velocity app
    app := velocity.New()

    // http middleware -> velocity.Handler
    app.Use(adaptor.HTTPMiddleware(logMiddleware))

    // Listen on port 3000
    app.Listen(":3000")
}

func logMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("log middleware")
        next.ServeHTTP(w, r)
    })
}
```

### Velocity Handler to net/http

```go
package main

import (
    "net/http"

    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/adaptor"
)

func main() {
    // velocity.Handler -> http.Handler
    http.Handle("/", adaptor.VelocityHandler(greet))

      // velocity.Handler -> http.HandlerFunc
    http.HandleFunc("/func", adaptor.VelocityHandlerFunc(greet))

    // Listen on port 3000
    http.ListenAndServe(":3000", nil)
}

func greet(c velocity.Ctx) error {
    return c.SendString("Hello World!")
}
```

### Velocity App to net/http

```go
package main

import (
    "net/http"

    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/adaptor"
)

func main() {
    app := velocity.New()

    app.Get("/greet", greet)

    // Listen on port 3000
    http.ListenAndServe(":3000", adaptor.VelocityApp(app))
}

func greet(c velocity.Ctx) error {
    return c.SendString("Hello World!")
}
```

### Velocity Context to (net/http).Request

```go
package main

import (
    "net/http"

    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/adaptor"
)

func main() {
    app := velocity.New()

    app.Get("/greet", greetWithHTTPReq)

    // Listen on port 3000
    http.ListenAndServe(":3000", adaptor.VelocityApp(app))
}

func greetWithHTTPReq(c velocity.Ctx) error {
    httpReq, err := adaptor.ConvertRequest(c, false)
    if err != nil {
        return err
    }

    return c.SendString("Request URL: " + httpReq.URL.String())
}
```
