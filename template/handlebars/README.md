---
id: handlebars
title: Handlebars
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=handlebars*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Handlebars is a template engine create by [aymerick](<<<https://github.com/aymerick/raymond), to >>>see the original syntax documentation please [click here](<<<https://github.com/aymerick/raymond#table-of-content>>>s)

### Basic Example

_**./views/index.hbs**_
```html  

{{> 'partials/header' }}

<h1>{{Title}}</h1>

{{> 'partials/footer' }}
```  

_**./views/partials/header.hbs**_
```html
<h2>Header</h2>
```  

_**./views/partials/footer.hbs**_
```html
<h2>Footer</h2>
```  

_**./views/layouts/main.hbs**_
```html
<!DOCTYPE html>
<html>

<head>
  <title>Main</title>
</head>

<body>
  {{embed}}
</body>

</html>
```  


```go
package main

import (
	"log"
	
	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/handlebars/v2"
)

func main() {
	// Create a new engine
	engine := handlebars.New("./views", ".hbs")

  // Or from an embedded system
  // See github.com/khulnasoft/embed for examples
  // engine := html.NewFileSystem(http.Dir("./views", ".hbs"))

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
