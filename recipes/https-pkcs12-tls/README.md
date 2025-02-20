---
title: HTTPS with PKCS12 TLS
keywords: [https, tls, pkcs12]
description: Setting up an HTTPS server with PKCS12 TLS certificates.
---

# HTTPS with PKCS12 TLS Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/https-pkcs12-tls) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/https-pkcs12-tls)

This project demonstrates how to set up an HTTPS server with PKCS12 TLS in a Go application using the Velocity framework.

## Prerequisites

Ensure you have the following installed:

- Golang
- [Velocity](https://github.com/khulnasoft/velocity) package
- PKCS12 certificate file (`cert.p12`)

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/https-pkcs12-tls
    ```

2. Install dependencies:
    ```sh
    go get
    ```

3. Place your PKCS12 certificate file (`cert.p12`) in the project directory.

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. Access the application at `https://localhost:3000`.

## Example

Here is an example of how to set up an HTTPS server with PKCS12 TLS in a Velocity application:

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
    "log"

    "go.khulnasoft.com/velocity"
    "golang.org/x/crypto/pkcs12"
)

func main() {
    // Load PKCS12 certificate
    p12Data, err := ioutil.ReadFile("cert.p12")
    if err != nil {
        log.Fatal(err)
    }

    // Decode PKCS12 certificate
    blocks, err := pkcs12.ToPEM(p12Data, "password")
    if err != nil {
        log.Fatal(err)
    }

    var pemData []byte
    for _, b := range blocks {
        pemData = append(pemData, pem.EncodeToMemory(b)...)
    }

    // Load certificate and key
    cert, err := tls.X509KeyPair(pemData, pemData)
    if err != nil {
        log.Fatal(err)
    }

    // Create TLS configuration
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    x509.NewCertPool(),
    }

    // Velocity instance
    app := velocity.New()

    // Routes
    app.Get("/", func(c *velocity.Ctx) error {
        return c.SendString("Hello, HTTPS with PKCS12 TLS!")
    })

    // Start server with TLS
    log.Fatal(app.ListenTLS(":3000", tlsConfig))
}
```

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [TLS in Go](https://golang.org/pkg/crypto/tls/)
- [PKCS12 in Go](https://pkg.go.dev/golang.org/x/crypto/pkcs12)
