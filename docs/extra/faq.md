---
id: faq
title: 🤔 FAQ
description: >-
  List of frequently asked questions. Feel free to open an issue to add your
  question to this page.
sidebar_position: 1
---

## How should I structure my application?

There is no definitive answer to this question. The answer depends on the scale of your application and the team that is involved. To be as flexible as possible, Velocity makes no assumptions in terms of structure.

Routes and other application-specific logic can live in as many files as you wish, in any directory structure you prefer. View the following examples for inspiration:

* [khulnasoft/boilerplate](https://github.com/khulnasoft/boilerplate)
* [thomasvvugt/velocity-boilerplate](https://github.com/thomasvvugt/velocity-boilerplate)
* [Youtube - Building a REST API using Gorm and Velocity](https://www.youtube.com/watch?v=Iq2qT0fRhAA)
* [embedmode/velocityseed](https://github.com/embedmode/velocityseed)

## How do I handle custom 404 responses?

If you're using v2.32.0 or later, all you need to do is to implement a custom error handler. See below, or see a more detailed explanation at [Error Handling](../guide/error-handling.md#custom-error-handler).

If you're using v2.31.0 or earlier, the error handler will not capture 404 errors. Instead, you need to add a middleware function at the very bottom of the stack \(below all other functions\) to handle a 404 response:

```go title="Example"
app.Use(func(c velocity.Ctx) error {
    return c.Status(velocity.StatusNotFound).SendString("Sorry can't find that!")
})
```

## How can i use live reload ?

[Air](https://github.com/air-verse/air) is a handy tool that automatically restarts your Go applications whenever the source code changes, making your development process faster and more efficient.

To use Air in a Velocity project, follow these steps:

* Install Air by downloading the appropriate binary for your operating system from the GitHub release page or by building the tool directly from source.
* Create a configuration file for Air in your project directory. This file can be named, for example, .air.toml or air.conf. Here's a sample configuration file that works with Velocity:

```toml
# .air.toml
root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000 # ms
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test\\.go"]
```

* Start your Velocity application using Air by running the following command in the terminal:

```sh
air
```

As you make changes to your source code, Air will detect them and automatically restart the application.

A complete example demonstrating the use of Air with Velocity can be found in the [Velocity Recipes repository](https://github.com/khulnasoft/recipes/tree/master/air). This example shows how to configure and use Air in a Velocity project to create an efficient development environment.

## How do I set up an error handler?

To override the default error handler, you can override the default when providing a [Config](../api/velocity.md#errorhandler) when initiating a new [Velocity instance](../api/velocity.md#new).

```go title="Example"
app := velocity.New(velocity.Config{
    ErrorHandler: func(c velocity.Ctx, err error) error {
        return c.Status(velocity.StatusInternalServerError).SendString(err.Error())
    },
})
```

We have a dedicated page explaining how error handling works in Velocity, see [Error Handling](../guide/error-handling.md).

## Which template engines does Velocity support?

Velocity currently supports 9 template engines in our [khulnasoft/template](https://docs.khulnasoft.io/template/) middleware:

* [ace](https://docs.khulnasoft.io/template/ace/)
* [amber](https://docs.khulnasoft.io/template/amber/)
* [django](https://docs.khulnasoft.io/template/django/)
* [handlebars](https://docs.khulnasoft.io/template/handlebars/)
* [html](https://docs.khulnasoft.io/template/html/)
* [jet](https://docs.khulnasoft.io/template/jet/)
* [mustache](https://docs.khulnasoft.io/template/mustache/)
* [pug](https://docs.khulnasoft.io/template/pug/)
* [slim](https://docs.khulnasoft.io/template/slim/)

To learn more about using Templates in Velocity, see [Templates](../guide/templates.md).

## Does Velocity have a community chat?

Yes, we have our own [Discord](https://khulnasoft.io/discord)server, where we hang out. We have different rooms for every subject.  
If you have questions or just want to have a chat, feel free to join us via this **&gt;** [**invite link**](https://khulnasoft.io/discord) **&lt;**.

![](/img/support-discord.png)

## Does velocity support sub domain routing ?

Yes we do, here are some examples:
This example works v2

```go
package main

import (
    "log"

    "go.khulnasoft.com/velocity/v3"
    "go.khulnasoft.com/velocity/v3/middleware/logger"
)

type Host struct {
    Velocity *velocity.App
}

func main() {
    // Hosts
    hosts := map[string]*Host{}
    //-----
    // API
    //-----
    api := velocity.New()
    api.Use(logger.New(logger.Config{
        Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
    }))
    hosts["api.localhost:3000"] = &Host{api}
    api.Get("/", func(c velocity.Ctx) error {
        return c.SendString("API")
    })
    //------
    // Blog
    //------
    blog := velocity.New()
    blog.Use(logger.New(logger.Config{
        Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
    }))
    hosts["blog.localhost:3000"] = &Host{blog}
    blog.Get("/", func(c velocity.Ctx) error {
        return c.SendString("Blog")
    })
    //---------
    // Website
    //---------
    site := velocity.New()
    site.Use(logger.New(logger.Config{
        Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
    }))

    hosts["localhost:3000"] = &Host{site}
    site.Get("/", func(c velocity.Ctx) error {
        return c.SendString("Website")
    })
    // Server
    app := velocity.New()
    app.Use(func(c velocity.Ctx) error {
        host := hosts[c.Hostname()]
        if host == nil {
            return c.SendStatus(velocity.StatusNotFound)
        } else {
            host.Velocity.Handler()(c.Context())
            return nil
        }
    })
    log.Fatal(app.Listen(":3000"))
}
```

If more information is needed, please refer to this issue [#750](https://go.khulnasoft.com/velocity/issues/750)
