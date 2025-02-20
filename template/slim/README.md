---
id: slim
title: Slim
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=slim*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Slim is a template engine created by [mattn](<<<https://github.com/mattn/go->>>slim), to see the original syntax documentation please [click here](<<<https://rubydoc.info/gem>>>s/slim/frames)

### Basic Example

_**./views/index.slim**_
```html  

== render("partials/header.slim")

h1 = Title

== render("partials/footer.slim")
```  

_**./views/partials/header.slim**_
```html
h2 = Header
```  

_**./views/partials/footer.slim**_
```html
h2 = Footer
```  

_**./views/layouts/main.slim**_
```html
doctype html
html
  head
    title Main
    include ../partials/meta.slim
  body
    == embed
```  


```go
package main

import (
    "log"

    "github.com/khulnasoft/velocity"
    "github.com/khulnasoft/velocity/template/slim/v2"

    // "net/http" // embedded system
)

func main() {
    // Create a new engine
    engine := slim.New("./views", ".slim")

    // Or from an embedded system
    // See github.com/khulnasoft/embed for examples
    // engine := slim.NewFileSystem(http.Dir("./views", ".slim"))

    // Pass the engine to the Views
    app := velocity.New(velocity.Config{
        Views: engine,
    })

    app.Get("/", func(c *velocity.Ctx) error {
        // Render index
        return c.Render("index", velocity.Map{
            "Title": "Hello, World!",
        })
    })

    app.Get("/layout", func(c *velocity.Ctx) error {
        // Render index within layouts/main
        return c.Render("index", velocity.Map{
            "Title": "Hello, World!",
        }, "layouts/main")
    })

    log.Fatal(app.Listen(":3000"))
}

```
