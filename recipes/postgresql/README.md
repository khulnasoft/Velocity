---
title: PostgreSQL
keywords: [postgresql]
description: Connecting to a PostgreSQL database.
---

# PostgreSQL Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/postgresql) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/postgresql)

This project demonstrates how to connect to a PostgreSQL database in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package
- PostgreSQL

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/postgresql
    ```

2. Install dependencies:
    ```sh
    go get
    ```

3. Set up your PostgreSQL database and update the connection string in the code.

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. Access the application at `http://localhost:3000`.

## Example

Here is an example of how to connect to a PostgreSQL database in a Velocity application:

```go
package main

import (
    "database/sql"
    "log"

    "go.khulnasoft.com/velocity"
    _ "github.com/lib/pq"
)

func main() {
    // Database connection
    connStr := "user=username dbname=mydb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Velocity instance
    app := velocity.New()

    // Routes
    app.Get("/", func(c *velocity.Ctx) error {
        var greeting string
        err := db.QueryRow("SELECT 'Hello, World!'").Scan(&greeting)
        if err != nil {
            return err
        }
        return c.SendString(greeting)
    })

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [pq Driver Documentation](https://pkg.go.dev/github.com/lib/pq)
