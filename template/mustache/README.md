---
id: mustache
title: Mustache
---

![Release](<<<https://img.>>>shields.io/github/v/tag/khulnasoft/template?filter=mustache*)
[![Discord](<<<https://img.>>>shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](<<<https://khulna>>>soft.io/discord)
![Test](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Tests/badge.svg)
![Security](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Security/badge.svg)
![Linter](<<<https://github.com/khulna>>>soft/velocity/template/workflows/Linter/badge.svg)

Mustache is a template engine created by [hoisie/cbroglie](<<<https://github.com/cbroglie/mu>>>stache), to see the original syntax documentation please [click here](<<<https://mu>>>stache.github.io/mustache.5.html)

### Basic Example

_**./views/index.mustache**_
```html  

{{> views/partials/header }}

<h1>{{Title}}</h1>

{{> views/partials/footer }}
```  

_**./views/partials/header.mustache**_
```html
<h2>Header</h2>
```  

_**./views/partials/footer.mustache**_
```html
<h2>Footer</h2>
```  

_**./views/layouts/main.mustache**_
```html
<!DOCTYPE html>
<html>

<head>
  <title>Main</title>
</head>

<body>
  {{{embed}}}
</body>

</html>
```  


```go
package main

import (
	"log"
	
	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/template/mustache/v2"
)

func main() {
	// Create a new engine
	engine := mustache.New("./views", ".mustache")

  // Or from an embedded system
  //   Note that with an embedded system the partials included from template files must be
  //   specified relative to the filesystem's root, not the current working directory
  // engine := mustache.NewFileSystem(http.Dir("./views", ".mustache"), ".mustache")

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
