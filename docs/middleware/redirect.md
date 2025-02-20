---
id: redirect
---

# Redirect

Redirection middleware for Velocity.

## Signatures

```go
func New(config ...Config) velocity.Handler
```

## Examples

```go
package main

import (
    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/middleware/redirect"
)

func main() {
    app := velocity.New()
    
    app.Use(redirect.New(redirect.Config{
      Rules: map[string]string{
        "/old":   "/new",
        "/old/*": "/new/$1",
      },
      StatusCode: 301,
    }))
    
    app.Get("/new", func(c velocity.Ctx) error {
      return c.SendString("Hello, World!")
    })
    app.Get("/new/*", func(c velocity.Ctx) error {
      return c.SendString("Wildcard: " + c.Params("*"))
    })
    
    app.Listen(":3000")
}
```

## Test

```bash
curl http://localhost:3000/old
curl http://localhost:3000/old/hello
```

## Config

| Property   | Type                    | Description                                                                                                                | Default                |
|:-----------|:------------------------|:---------------------------------------------------------------------------------------------------------------------------|:-----------------------|
| Next       | `func(velocity.Ctx) bool` | Filter defines a function to skip middleware.                                                                              | `nil`                  |
| Rules      | `map[string]string`     | Rules defines the URL path rewrite rules. The values captured in asterisk can be retrieved by index e.g. $1, $2 and so on. | Required               |
| StatusCode | `int`                   | The status code when redirecting. This is ignored if Redirect is disabled.                                                 | 302 Temporary Redirect |

## Default Config

```go
var ConfigDefault = Config{
    StatusCode: velocity.StatusFound,
}
```
