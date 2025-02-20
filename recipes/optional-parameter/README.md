---
title: Optional Parameter
keywords: [optional, parameter]
description: Handling optional parameters.
---

# Optional Parameter Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/optional-parameter) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/optional-parameter)

This project demonstrates how to handle optional parameters in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/optional-parameter
    ```

2. Install dependencies:
    ```sh
    go get
    ```

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

## Example

Here is an example of how to handle optional parameters in a Velocity application:

```go
package main

import (
    "go.khulnasoft.com/velocity"
)

func main() {
    app := velocity.New()

    app.Get("/user/:id?", func(c *velocity.Ctx) error {
        id := c.Params("id", "defaultID")
        return c.SendString("User ID: " + id)
    })

    app.Listen(":3000")
}
```

In this example:
- The `:id?` parameter in the route is optional.
- If the `id` parameter is not provided, it defaults to `"defaultID"`.

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
