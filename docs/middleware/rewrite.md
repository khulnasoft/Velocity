---
id: rewrite
---

# Rewrite

Rewrite middleware rewrites the URL path based on provided rules. It can be helpful for backward compatibility or just creating cleaner and more descriptive links.

## Signatures

```go
func New(config ...Config) velocity.Handler
```

## Config

| Property | Type                    | Description                                                                                          | Default    |
|:---------|:------------------------|:-----------------------------------------------------------------------------------------------------|:-----------|
| Next     | `func(velocity.Ctx) bool` | Next defines a function to skip middleware.                                                          | `nil`      |
| Rules    | `map[string]string`     | Rules defines the URL path rewrite rules. The values captured in asterisk can be retrieved by index. | (Required) |

### Examples

```go
package main

import (
    "go.khulnasoft.com/velocity"
    "go.khulnasoft.com/velocity/middleware/rewrite"
)

func main() {
    app := velocity.New()
    
    app.Use(rewrite.New(rewrite.Config{
      Rules: map[string]string{
        "/old":   "/new",
        "/old/*": "/new/$1",
      },
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
