---
id: ace
title: Ace
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=ace*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Ace is a template engine create by [yossi](<<<https://github.com/yo>>>sssi/ace), to see the original syntax documentation please [click here](<<<https://github.com/yo>>>sssi/ace/blob/master/documentation/syntax.md)

### Basic Example

_**./views/index.ace**_
```html  

= include ./views/partials/header .

h1 {{.Title}}

= include ./views/partials/footer .
```  

_**./views/partials/header.ace**_
```html
h1 Header
```  

_**./views/partials/footer.ace**_
```html
h1 Footer
```  

_**./views/layouts/main.ace**_
```html
= doctype html
html
  head
    title Main
  body
    {{embed}}
```  


```go
package main

import (
	"log"
	
	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/ace/v2"
)

func main() {
	// Create a new engine
	engine := ace.New("./views", ".ace")

  // Or from an embedded system
  // See github.com/khulnasoft/embed for examples
  // engine := html.NewFileSystem(http.Dir("./views", ".ace"))

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
