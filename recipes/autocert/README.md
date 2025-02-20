---
title: Autocert
keywords: [autocert, tls, letsencrypt, ssl, https, certificate]
description: Automatic TLS certificate management.
---

# Autocert Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/autocert) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/autocert)

This example demonstrates how to set up a secure Go Velocity application using Let's Encrypt for automatic TLS certificate management with `autocert`.

## Description

This project provides a starting point for building a secure web application with automatic TLS certificate management using Let's Encrypt. It leverages Velocity for the web framework and `autocert` for certificate management.

## Requirements

- [Go](https://golang.org/dl/) 1.18 or higher
- [Git](https://git-scm.com/downloads)

## Setup

1. Clone the repository:
    ```bash
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/autocert
    ```

2. Install the dependencies:
    ```bash
    go mod download
    ```

3. Update the `HostPolicy` in `main.go` with your domain:
    ```go
    m := &autocert.Manager{
        Prompt: autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist("yourdomain.com"), // Replace with your domain
        Cache: autocert.DirCache("./certs"),
    }
    ```

4. Run the application:
    ```bash
    go run main.go
    ```

The application should now be running on `https://localhost`.

## Example Usage

1. Open your browser and navigate to `https://yourdomain.com` (replace with your actual domain).

2. You should see the message: `This is a secure server 👮`.

## Conclusion

This example provides a basic setup for a Go Velocity application with automatic TLS certificate management using Let's Encrypt. It can be extended and customized further to fit the needs of more complex applications.

## References

- [Velocity Documentation](https://docs.khulnasoft.io)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Autocert Documentation](https://pkg.go.dev/golang.org/x/crypto/acme/autocert)
