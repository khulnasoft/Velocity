---
id: jet
title: Jet
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=jet*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Jet is a template engine create by [cloudykit](<<<https://github.com/CloudyKit/jet), to >>>see the original syntax documentation please [click here](<<<https://github.com/CloudyKit/jet/wiki/3.-Jet-template->>>syntax)

### Basic Example

_**./views/index.jet**_
```html  

{{include "partials/header"}}

<h1>{{ Title }}</h1>

{{include "partials/footer"}}
```  

_**./views/partials/header.jet**_
```html
<h2>Header</h2>
```  

_**./views/partials/footer.jet**_
```html
<h2>Footer</h2>
```  

_**./views/layouts/main.jet**_
```html
<!DOCTYPE html>
<html>

<head>
  <title>Title</title>
</head>

<body>
  {{ embed() }}
</body>

</html>
```  


```go
package main

import (
	"log"
	
	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/jet/v2"
)

func main() {
	// Create a new engine
	engine := jet.New("./views", ".jet")

	// Or from an embedded system
	// See github.com/khulnasoft/embed for examples
	// engine := jet.NewFileSystem(http.Dir("./views", ".jet"))

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
