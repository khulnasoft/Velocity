---
id: timeout
---

# Timeout

There exist two distinct implementations of timeout middleware [Velocity](https://github.com/khulnasoft/velocity).

## New

As a `velocity.Handler` wrapper, it creates a context with `context.WithTimeout` which is then used with `c.Context()`.

If the context passed executions (eg. DB ops, Http calls) takes longer than the given duration to return, the timeout error is set and forwarded to the centralized `ErrorHandler`.

It does not cancel long running executions. Underlying executions must handle timeout by using `context.Context` parameter.

## Signatures

```go
func New(handler velocity.Handler, timeout time.Duration, timeoutErrors ...error) velocity.Handler
```

## Examples

Import the middleware package that is part of the Velocity web framework

```go
import (
    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/middleware/timeout"
)
```

After you initiate your Velocity app, you can use the following possibilities:

```go
func main() {
    app := velocity.New()
    h := func(c velocity.Ctx) error {
        sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
        if err := sleepWithContext(c.Context(), sleepTime); err != nil {
            return fmt.Errorf("%w: execution error", err)
        }
        return nil
    }

    app.Get("/foo/:sleepTime", timeout.New(h, 2*time.Second))
    log.Fatal(app.Listen(":3000"))
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
    timer := time.NewTimer(d)

    select {
    case <-ctx.Done():
        if !timer.Stop() {
            <-timer.C
        }
        return context.DeadlineExceeded
    case <-timer.C:
    }
    return nil
}
```

Test http 200 with curl:

```bash
curl --location -I --request GET 'http://localhost:3000/foo/1000' 
```

Test http 408 with curl:

```bash
curl --location -I --request GET 'http://localhost:3000/foo/3000' 
```

Use with custom error:

```go
var ErrFooTimeOut = errors.New("foo context canceled")

func main() {
    app := velocity.New()
    h := func(c velocity.Ctx) error {
        sleepTime, _ := time.ParseDuration(c.Params("sleepTime") + "ms")
        if err := sleepWithContextWithCustomError(c.Context(), sleepTime); err != nil {
            return fmt.Errorf("%w: execution error", err)
        }
        return nil
    }

    app.Get("/foo/:sleepTime", timeout.New(h, 2*time.Second, ErrFooTimeOut))
    log.Fatal(app.Listen(":3000"))
}

func sleepWithContextWithCustomError(ctx context.Context, d time.Duration) error {
    timer := time.NewTimer(d)
    select {
    case <-ctx.Done():
        if !timer.Stop() {
            <-timer.C
        }
        return ErrFooTimeOut
    case <-timer.C:
    }
    return nil
}
```

Sample usage with a DB call:

```go
func main() {
    app := velocity.New()
    db, _ := gorm.Open(postgres.Open("postgres://localhost/foodb"), &gorm.Config{})

    handler := func(ctx velocity.Ctx) error {
        tran := db.WithContext(ctx.Context()).Begin()
        
        if tran = tran.Exec("SELECT pg_sleep(50)"); tran.Error != nil {
            return tran.Error
        }
        
        if tran = tran.Commit(); tran.Error != nil {
            return tran.Error
        }

        return nil
    }

    app.Get("/foo", timeout.New(handler, 10*time.Second))
    log.Fatal(app.Listen(":3000"))
}
```
