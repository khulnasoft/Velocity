---
title: File Server
keywords: [file server, static files]
description: Serving static files.
---

# File Server Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/file-server) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/file-server)

This project demonstrates how to set up a simple file server in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/file-server
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

2. Access the file server at `http://localhost:3000`.

## Example

Here is an example `main.go` file for the Velocity application serving static files:

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
)

func main() {
    app := velocity.New()

    // Serve static files from the "public" directory
    app.Static("/", "./public")

    log.Fatal(app.Listen(":3000"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Golang Documentation](https://golang.org/doc/)
