---
title: HTTPS with TLS
keywords: [https, tls, ssl, self-signed]
description: Setting up an HTTPS server with self-signed TLS certificates.
---

# HTTPS with TLS Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/https-tls) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/https-tls)

This project demonstrates how to set up an HTTPS server with TLS in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package
- TLS certificates (self-signed or from a trusted CA)

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/https-tls
    ```

2. Install dependencies:
    ```sh
    go get
    ```

3. Place your TLS certificate (`cert.pem`) and private key (`key.pem`) in the project directory.

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. Access the application at `https://localhost:3000`.

## Example

Here is an example of how to set up an HTTPS server with TLS in a Velocity application:

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity"
)

func main() {
    app := velocity.New()

    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, HTTPS with TLS!")
    })

    // Start server with TLS
    log.Fatal(app.ListenTLS(":3000", "cert.pem", "key.pem"))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [TLS in Go](https://golang.org/pkg/crypto/tls/)
