---
title: GraphQL
keywords: [graphql]
description: Setting up a GraphQL server.
---

# GraphQL Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/graphql) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/graphql)

This project demonstrates how to set up a GraphQL server in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package
- [gqlgen](https://github.com/99designs/gqlgen) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/graphql
    ```

2. Install dependencies:
    ```sh
    go get
    ```

3. Initialize gqlgen:
    ```sh
    go run github.com/99designs/gqlgen init
    ```

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. Access the GraphQL playground at `http://localhost:3000/graphql`.

## Example

Here is an example `main.go` file for the Velocity application with GraphQL:

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
)

func main() {
    app := velocity.New()

    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver{}}))

    app.All("/graphql", func(c *velocity.Ctx) error {
        srv.ServeHTTP(c.Context().ResponseWriter(), c.Context().Request)
        return nil
    })

    app.Get("/", func(c *velocity.Ctx) error {
        playground.Handler("GraphQL playground", "/graphql").ServeHTTP(c.Context().ResponseWriter(), c.Context().Request)
        return nil
    })

    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Documentation](https://graphql.org/)
