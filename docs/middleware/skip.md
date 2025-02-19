---
id: skip
---

# Skip

Skip middleware for [Velocity](https://go.khulnasoft.com/velocity) that skips a wrapped handler if a predicate is true.

## Signatures

```go
func New(handler velocity.Handler, exclude func(c velocity.Ctx) bool) velocity.Handler
```

## Examples

Import the middleware package that is part of the Velocity web framework

```go
import (
    "go.khulnasoft.com/velocity/v3"
    "go.khulnasoft.com/velocity/v3/middleware/skip"
)
```

After you initiate your Velocity app, you can use the following possibilities:

```go
func main() {
    app := velocity.New()

    app.Use(skip.New(BasicHandler, func(ctx velocity.Ctx) bool {
        return ctx.Method() == velocity.MethodGet
    }))

    app.Get("/", func(ctx velocity.Ctx) error {
        return ctx.SendString("It was a GET request!")
    })

    log.Fatal(app.Listen(":3000"))
}

func BasicHandler(ctx velocity.Ctx) error {
    return ctx.SendString("It was not a GET request!")
}
```

:::tip
app.Use will handle requests from any route, and any method. In the example above, it will only skip if the method is GET.
:::
