---
title: JWT
keywords: [jwt, json web token, authentication]
description: Using JSON Web Tokens (JWT) for authentication.
---

# Velocity with JWT

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/jwt) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/jwt)

This example demonstrates how to use JSON Web Tokens (JWT) for authentication in a Velocity application.

## Prerequisites

- Go 1.16 or higher
- Go modules

## Setup

1. Clone the repository:
    ```sh
    git clone https://go.khulnasoft.com/velocity/recipes.git
    cd recipes/jwt
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Running the Application

1. Run the application:
    ```sh
    go run main.go
    ```

2. The server will start on `http://localhost:3000`.

## Endpoints

| Method | URL           | Description                |
| ------ | ------------- | -------------------------- |
| POST   | /login        | Authenticates a user and returns a JWT |
| GET    | /restricted   | Accesses a restricted route with JWT   |

## Example Requests

### Login
```sh
curl -X POST http://localhost:3000/login -d '{"username": "user", "password": "pass"}' -H "Content-Type: application/json"
```

### Access Restricted Route
```sh
curl -X GET http://localhost:3000/restricted -H "Authorization: Bearer <your_jwt_token>"
```

## Postman Collection

You can find Postman examples [here](https://www.getpostman.com/collections/0e83876e0f2a0c8ecd70).
