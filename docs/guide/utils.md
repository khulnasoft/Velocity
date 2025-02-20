---
id: utils
title: 🧰 Utils
sidebar_position: 8
toc_max_heading_level: 4
---

## Generics

### Convert

Converts a string value to a specified type, handling errors and optional default values.
This function simplifies the conversion process by encapsulating error handling and the management of default values, making your code cleaner and more consistent.

```go title="Signature"
func Convert[T any](value string, convertor func(string) (T, error), defaultValue ...T) (*T, error)
```

```go title="Example"
// GET http://example.com/id/bb70ab33-d455-4a03-8d78-d3c1dacae9ff
app.Get("/id/:id", func(c velocity.Ctx) error {
   velocity.Convert(c.Params("id"), uuid.Parse)                   // UUID(bb70ab33-d455-4a03-8d78-d3c1dacae9ff), nil


// GET http://example.com/search?id=65f6f54221fb90e6a6b76db7
app.Get("/search", func(c velocity.Ctx) error) {
    velocity.Convert(c.Query("id"), mongo.ParseObjectID)           // objectid(65f6f54221fb90e6a6b76db7), nil
    velocity.Convert(c.Query("id"), uuid.Parse)                    // uuid.Nil, error(cannot parse given uuid)
    velocity.Convert(c.Query("id"), uuid.Parse, mongo.NewObjectID) // new object id generated and return nil as error.
}

  // ...
})
```

### GetReqHeader

GetReqHeader function utilizing Go's generics feature.
This function allows for retrieving HTTP request headers with a more specific data type.

```go title="Signature"
func GetReqHeader[V any](c Ctx, key string, defaultValue ...V) V
```

```go title="Example"
app.Get("/search", func(c velocity.Ctx) error {
    // curl -X GET http://example.com/search -H "X-Request-ID: 12345" -H "X-Request-Name: John"
    velocity.GetReqHeader[int](c, "X-Request-ID")               // => returns 12345 as integer.
    velocity.GetReqHeader[string](c, "X-Request-Name")          // => returns "John" as string.
    velocity.GetReqHeader[string](c, "unknownParam", "default") // => returns "default" as string.
    // ...
})
```

### Locals

Locals function utilizing Go's generics feature.
This function allows for manipulating and retrieving local values within a request context with a more specific data type.

```go title="Signature"
func Locals[V any](c Ctx, key any, value ...V) V

// get local value
func Locals[V any](c Ctx, key any) V
// set local value
func Locals[V any](c Ctx, key any, value ...V) V
```

```go title="Example"
app.Use("/user/:user/:id", func(c velocity.Ctx) error {
    // set local values
    velocity.Locals[string](c, "user", "john")
    velocity.Locals[int](c, "id", 25)
    // ...
    
    return c.Next()
})


app.Get("/user/*", func(c velocity.Ctx) error {
    // get local values
    name := velocity.Locals[string](c, "user") // john
    age := velocity.Locals[int](c, "id")       // 25
    // ...
})
```

### Params

Params function utilizing Go's generics feature.
This function allows for retrieving route parameters with a more specific data type.

```go title="Signature"
func Params[V any](c Ctx, key string, defaultValue ...V) V
```

```go title="Example"
app.Get("/user/:user/:id", func(c velocity.Ctx) error {
    // http://example.com/user/john/25
    velocity.Params[int](c, "id")               // => returns 25 as integer.
    velocity.Params[int](c, "unknownParam", 99) // => returns the default 99 as integer.
    // ...
    return c.SendString("Hello, " + velocity.Params[string](c, "user"))
})
```

### Query

Query function utilizing Go's generics feature.
This function allows for retrieving query parameters with a more specific data type.

```go title="Signature"
func Query[V any](c Ctx, key string, defaultValue ...V) V
```

```go title="Example"
app.Get("/search", func(c velocity.Ctx) error {
    // http://example.com/search?name=john&age=25
    velocity.Query[string](c, "name")                    // => returns "john"
    velocity.Query[int](c, "age")                        // => returns 25 as integer.
    velocity.Query[string](c, "unknownParam", "default") // => returns "default" as string.
    // ...
})
```
