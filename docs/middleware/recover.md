---
id: recover
---

# Recover

Recover middleware for [Velocity](https://go.khulnasoft.com/velocity) that recovers from panics anywhere in the stack chain and handles the control to the centralized [ErrorHandler](https://docs.khulnasoft.io/guide/error-handling).

## Signatures

```go
func New(config ...Config) velocity.Handler
```

## Examples

Import the middleware package that is part of the Velocity web framework

```go
import (
    "go.khulnasoft.com/velocity/v3"
    recoverer "go.khulnasoft.com/velocity/v3/middleware/recover"
)
```

After you initiate your Velocity app, you can use the following possibilities:

```go
// Initialize default config
app.Use(recoverer.New())

// This panic will be caught by the middleware
app.Get("/", func(c velocity.Ctx) error {
    panic("I'm an error")
})
```

## Config

| Property          | Type                            | Description                                                         | Default                  |
|:------------------|:--------------------------------|:--------------------------------------------------------------------|:-------------------------|
| Next              | `func(velocity.Ctx) bool`         | Next defines a function to skip this middleware when returned true. | `nil`                    |
| EnableStackTrace  | `bool`                          | EnableStackTrace enables handling stack trace.                      | `false`                  |
| StackTraceHandler | `func(velocity.Ctx, any)` | StackTraceHandler defines a function to handle stack trace.         | defaultStackTraceHandler |

## Default Config

```go
var ConfigDefault = Config{
    Next:              nil,
    EnableStackTrace:  false,
    StackTraceHandler: defaultStackTraceHandler,
}
```
