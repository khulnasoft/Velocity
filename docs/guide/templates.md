---
id: templates
title: 📝 Templates
description: Velocity supports server-side template engines.
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Templates are a great tool to render dynamic content without using a separate frontend framework.

## Template Engines

Velocity allows you to provide a custom template engine at app initialization.

```go
app := velocity.New(velocity.Config{
    // Pass in Views Template Engine
    Views: engine,

    // Default global path to search for views (can be overriden when calling Render())
    ViewsLayout: "layouts/main",

    // Enables/Disables access to `ctx.Locals()` entries in rendered views
    // (defaults to false)
    PassLocalsToViews: false,
})
```

### Supported Engines

The Velocity team maintains a [templates](https://docs.khulnasoft.com/template) package that provides wrappers for multiple template engines:

* [ace](https://docs.khulnasoft.com/template/ace/)
* [amber](https://docs.khulnasoft.com/template/amber/)
* [django](https://docs.khulnasoft.com/template/django/)
* [handlebars](https://docs.khulnasoft.com/template/handlebars)
* [html](https://docs.khulnasoft.com/template/html)
* [jet](https://docs.khulnasoft.com/template/jet)
* [mustache](https://docs.khulnasoft.com/template/mustache)
* [pug](https://docs.khulnasoft.com/template/pug)
* [slim](https://docs.khulnasoft.com/template/slim)

:::info
Custom template engines can implement the `Views` interface to be supported in Velocity.
:::

```go title="Views interface"
type Views interface {
    // Velocity executes Load() on app initialization to load/parse the templates
    Load() error

    // Outputs a template to the provided buffer using the provided template,
    // template name, and binded data
    Render(io.Writer, string, interface{}, ...string) error
}
```

:::note
The `Render` method is linked to the [**ctx.Render\(\)**](../api/ctx.md#render) function that accepts a template name and binding data.
:::

## Rendering Templates

Once an engine is set up, a route handler can call the [**ctx.Render\(\)**](../api/ctx.md#render) function with a template name and binded data to send the rendered template.

```go title="Signature"
func (c Ctx) Render(name string, bind Map, layouts ...string) error
```

:::info
By default, [**ctx.Render\(\)**](../api/ctx.md#render) searches for the template name in the `ViewsLayout` path. To override this setting, provide the path(s) in the `layouts` argument.
:::

<Tabs>
<TabItem value="example" label="Example">

```go
app.Get("/", func(c velocity.Ctx) error {
    return c.Render("index", velocity.Map{
        "Title": "Hello, World!",
    })

})
```

</TabItem>

<TabItem value="index" label="layouts/index.html">

```html
<!DOCTYPE html>
<html>
    <body>
        <h1>{{.Title}}</h1>
    </body>
</html>
```

</TabItem>

</Tabs>

:::caution
If the Velocity config option `PassLocalsToViews` is enabled, then all locals set using `ctx.Locals(key, value)` will be passed to the template. It is important to avoid clashing keys when using this setting.
:::

## Advanced Templating

### Custom Functions

Velocity supports adding custom functions to templates.

#### AddFunc

Adds a global function to all templates.

```go title="Signature"
func (e *Engine) AddFunc(name string, fn interface{}) IEngineCore
```

<Tabs>
<TabItem value="add-func-example" label="AddFunc Example">

```go
// Add `ToUpper` to engine
engine := html.New("./views", ".html")
engine.AddFunc("ToUpper", func(s string) string {
    return strings.ToUpper(s)
}

// Initialize Velocity App
app := velocity.New(velocity.Config{
    Views: engine,
})

app.Get("/", func (c velocity.Ctx) error {
    return c.Render("index", velocity.Map{
        "Content": "hello, World!"
    })
})
```

</TabItem>
<TabItem value="add-func-template" label="views/index.html">

```html
<!DOCTYPE html>
<html>
    <body>
        <p>This will be in {{ToUpper "all caps"}}:</p>
        <p>{{ToUpper .Content}}</p>
    </body>
</html>
```

</TabItem>
</Tabs>

#### AddFuncMap

Adds a Map of functions (keyed by name) to all templates.

```go title="Signature"
func (e *Engine) AddFuncMap(m map[string]interface{}) IEngineCore
```

<Tabs>
<TabItem value="add-func-map-example" label="AddFuncMap Example">

```go
// Add `ToUpper` to engine
engine := html.New("./views", ".html")
engine.AddFuncMap(map[string]interface{}{
    "ToUpper": func(s string) string {
        return strings.ToUpper(s)
    },
})

// Initialize Velocity App
app := velocity.New(velocity.Config{
    Views: engine,
})

app.Get("/", func (c velocity.Ctx) error {
    return c.Render("index", velocity.Map{
        "Content": "hello, world!"
    })
})
```

</TabItem>
<TabItem value="add-func-map-template" label="views/index.html">

```html
<!DOCTYPE html>
<html>
    <body>
        <p>This will be in {{ToUpper "all caps"}}:</p>
        <p>{{ToUpper .Content}}</p>
    </body>
</html>
```

</TabItem>
</Tabs>

* For more advanced template documentation, please visit the [khulnasoft/template GitHub Repository](https://github.com/khulnasoft/template).

## Full Example

<Tabs>
<TabItem value="example" label="Example">

```go
package main

import (
    "log"
    "go.khulnasoft.com/velocity/v3"
    "github.com/khulnasoft/template/html/v2"
)

func main() {
    // Initialize standard Go html template engine
    engine := html.New("./views", ".html")
    // If you want to use another engine,
    // just replace with following:
    // Create a new engine with django
    // engine := django.New("./views", ".django")

    app := velocity.New(velocity.Config{
        Views: engine,
    })
    app.Get("/", func(c velocity.Ctx) error {
        // Render index template
        return c.Render("index", velocity.Map{
            "Title": "Go Velocity Template Example",
            "Description": "An example template",
            "Greeting": "Hello, World!",
        });
    })

    log.Fatal(app.Listen(":3000"))
}
```

</TabItem>
<TabItem value="index" label="views/index.html">

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{.Title}}</title>
        <meta name="description" content="{{.Description}}">
    </head>
<body>
    <h1>{{.Title}}</h1>
        <p>{{.Greeting}}</p>
</body>
</html>
```

</TabItem>
</Tabs>
