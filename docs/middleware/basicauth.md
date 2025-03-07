---
id: basicauth
---

# BasicAuth

Basic Authentication middleware for [Velocity](https://github.com/khulnasoft/velocity) that provides an HTTP basic authentication. It calls the next handler for valid credentials and [401 Unauthorized](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401) or a custom response for missing or invalid credentials.

## Signatures

```go
func New(config Config) velocity.Handler
func UsernameFromContext(c velocity.Ctx) string
func PasswordFromContext(c velocity.Ctx) string
```

## Examples

Import the middleware package that is part of the Velocity web framework

```go
import (
    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/middleware/basicauth"
)
```

After you initiate your Velocity app, you can use the following possibilities:

```go
// Provide a minimal config
app.Use(basicauth.New(basicauth.Config{
    Users: map[string]string{
        "john":  "doe",
        "admin": "123456",
    },
}))

// Or extend your config for customization
app.Use(basicauth.New(basicauth.Config{
    Users: map[string]string{
        "john":  "doe",
        "admin": "123456",
    },
    Realm: "Forbidden",
    Authorizer: func(user, pass string) bool {
        if user == "john" && pass == "doe" {
            return true
        }
        if user == "admin" && pass == "123456" {
            return true
        }
        return false
    },
    Unauthorized: func(c velocity.Ctx) error {
        return c.SendFile("./unauthorized.html")
    },
}))
```

Getting the username and password

```go
func handler(c velocity.Ctx) error {
    username := basicauth.UsernameFromContext(c)
    password := basicauth.PasswordFromContext(c)
    log.Printf("Username: %s Password: %s", username, password)
    return c.SendString("Hello, " + username)
}
```

## Config

| Property        | Type                        | Description                                                                                                                                                           | Default               |
|:----------------|:----------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------|:----------------------|
| Next            | `func(velocity.Ctx) bool`     | Next defines a function to skip this middleware when returned true.                                                                                                   | `nil`                 |
| Users           | `map[string]string`         | Users defines the allowed credentials.                                                                                                                                | `map[string]string{}` |
| Realm           | `string`                    | Realm is a string to define the realm attribute of BasicAuth. The realm identifies the system to authenticate against and can be used by clients to save credentials. | `"Restricted"`        |
| Authorizer      | `func(string, string) bool` | Authorizer defines a function to check the credentials. It will be called with a username and password and is expected to return true or false to indicate approval.  | `nil`                 |
| Unauthorized    | `velocity.Handler`             | Unauthorized defines the response body for unauthorized responses.                                                                                                    | `nil`                 |

## Default Config

```go
var ConfigDefault = Config{
    Next:            nil,
    Users:           map[string]string{},
    Realm:           "Restricted",
    Authorizer:      nil,
    Unauthorized:    nil,
}
```
