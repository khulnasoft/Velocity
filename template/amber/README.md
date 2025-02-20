---
id: amber
title: Amber
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=amber*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Amber is a template engine create by [eknkc](<<<https://github.com/eknkc/amber), to >>>see the original syntax documentation please [click here](<<<https://github.com/eknkc/amber#tag>>>s)

### Basic Example

_**./views/index.amber**_
```html  

import ./views/partials/header

h1 #{Title}

import ./views/partials/footer
```  

_**./views/partials/header.amber**_
```html
h1 Header
```  

_**./views/partials/footer.amber**_
```html
h1 Footer
```  

_**./views/layouts/main.amber**_
```html
doctype html
html
  head
    title Main
  body
    #{embed()}
```  


```go
package main

import (
	"log"
	
	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/amber/v2"
)

func main() {
	// Create a new engine
	engine := amber.New("./views", ".amber")

  // Or from an embedded system
  // See github.com/khulnasoft/embed for examples
  // engine := html.NewFileSystem(http.Dir("./views", ".amber"))

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
