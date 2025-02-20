---
title: Bootstrap
keywords: [bootstrap, gorm, validator, env]
description: Integrating Bootstrap.
---

# Bootstrap

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://go.khulnasoft.com/velocity/recipes/tree/master/bootstrap) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/khulnasoft/recipes/tree/master/bootstrap)

Velocity bootstrap for rapid development using Go-Velocity / Gorm / Validator.

## Components
* Velocity
  * Html Engine Template
  * Logger
  * Monitoring
* Gorm
  * PGSQL Driver
* Validator
* Env File

## Router
API Router `/api` with rate limiter middleware
Http Router `/` with CORS and CSRF middleware

## Setup

1. Copy the example env file over:
    ```
    cp .env.example .env
    ```

2. Modify the env file you just copied `.env` with the correct credentials for your database. Make sure the database you entered in `DB_NAME` has been created.

3. Run the API:
    ```
    go run main.go
    ```
Your api should be running at `http://localhost:4000/` if the port is in use you may modify it in the `.env` you just created.
