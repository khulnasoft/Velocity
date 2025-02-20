---
id: pug
title: Pug
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=pug*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Pug is a template engine create by [joker](<<<https://github.com/Joker/jade), to >>>see the original syntax documentation please [click here](<<<https://pugj>>>s.org/language/tags.html)

### Basic Example

_**./views/index.pug**_
```html  

include partials/header.pug

h1 #{.Title}

include partials/footer.pug
```  

_**./views/partials/header.pug**_
```html
h2 Header
```  

_**./views/partials/footer.pug**_
```html
h2 Footer
```  

_**./views/layouts/main.pug**_
```html
doctype html
html
  head
    title Main
    include ../partials/meta.pug
  body
    | {{embed}}
```  


```go
package main

import (
	"log"

	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/pug/v2"

	// "net/http" // embedded system
)

func main() {
	// Create a new engine
	engine := pug.New("./views", ".pug")

	// Or from an embedded system
	// See github.com/khulnasoft/embed for examples
	// engine := pug.NewFileSystem(http.Dir("./views"), ".pug")

	// Pass the engine to the views
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
