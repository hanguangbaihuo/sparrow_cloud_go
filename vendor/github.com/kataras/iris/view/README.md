# View

Iris supports 8 template engines out-of-the-box, developers can still use any external golang template engine,
as `Context.ResponseWriter()` is an `io.Writer`.

All template engines share a common API i.e.
Parse using embedded assets, Layouts and Party-specific layout, Template Funcs, Partial Render and more.

| #  | Name       | Parser   |
|:---|:-----------|----------|
| 1 | HTML       | [html/template](https://pkg.go.dev/html/template) |
| 2 | Blocks     | [kataras/blocks](https://github.com/kataras/blocks) |
| 3 | Django     | [flosch/pongo2](https://github.com/flosch/pongo2) |
| 4 | Pug        | [Joker/jade](https://github.com/Joker/jade) |
| 5 | Handlebars | [aymerick/raymond](https://github.com/aymerick/raymond) |
| 6 | Amber      | [eknkc/amber](https://github.com/eknkc/amber) |
| 7 | Jet        | [CloudyKit/jet](https://github.com/CloudyKit/jet) |
| 8 | Ace        | [yosssi/ace](https://github.com/yosssi/ace) |

[List of Examples](https://github.com/kataras/iris/tree/master/_examples/view).

You can serve [quicktemplate](https://github.com/valyala/quicktemplate) files too, simply by using the `Context.ResponseWriter`, take a look at the [iris/_examples/view/quicktemplate](https://github.com/kataras/iris/tree/master/_examples/view/quicktemplate) example.

## Overview

```go
// file: main.go
package main

import "github.com/kataras/iris"

func main() {
    app := iris.New()
    // Load all templates from the "./views" folder
    // where extension is ".html" and parse them
    // using the standard `html/template` package.
    app.RegisterView(iris.HTML("./views", ".html"))

    // Method:    GET
    // Resource:  http://localhost:8080
    app.Get("/", func(ctx iris.Context) {
        // Bind: {{.message}} with "Hello world!"
        ctx.ViewData("message", "Hello world!")
        // Render template file: ./views/hello.html
        ctx.View("hello.html")
    })

    // Method:    GET
    // Resource:  http://localhost:8080/user/42
    app.Get("/user/{id:int64}", func(ctx iris.Context) {
        userID, _ := ctx.Params().GetInt64("id")
        ctx.Writef("User ID: %d", userID)
    })

    // Start the server using a network address.
    app.Listen(":8080")
}
```

```html
<!-- file: ./views/hello.html -->
<html>
<head>
    <title>Hello Page</title>
</head>
<body>
    <h1>{{.message}}</h1>
</body>
</html>
```

## Template functions

```go
package main

import "github.com/kataras/iris"

func main() {
    app := iris.New()
    tmpl := iris.HTML("./templates", ".html")

    // builtin template funcs are:
    //
    // - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
    // - {{ render "header.html" }}
    // - {{ render_r "header.html" }} // partial relative path to current page
    // - {{ yield }}
    // - {{ current }}

    // register a custom template func.
    tmpl.AddFunc("greet", func(s string) string {
        return "Greetings " + s + "!"
    })

    // register the view engine to the views, this will load the templates.
    app.RegisterView(tmpl)

    app.Get("/", hi)

    // http://localhost:8080
    app.Listen(":8080")
}

func hi(ctx iris.Context) {
    // render the template file "./templates/hi.html"
    ctx.View("hi.html")
}
```

```html
<!-- file: ./templates/hi.html -->
<b>{{greet "kataras"}}</b> <!-- will be rendered as: <b>Greetings kataras!</b> -->
```

## Embedded

View engine supports bundled(https://github.com/go-bindata/go-bindata) template files too.
`go-bindata` gives you two functions, `Assset` and `AssetNames`,
these can be set to each of the template engines using the `.Binary` function.

Example code:

```go
package main

import "github.com/kataras/iris"

func main() {
    app := iris.New()
    // $ go get -u github.com/go-bindata/go-bindata/v3/go-bindata
    // $ go-bindata ./templates/...
    // $ go build
    // $ ./embedding-templates-into-app
    // html files are not used, you can delete the folder and run the example
    app.RegisterView(iris.HTML("./templates", ".html").Binary(Asset, AssetNames))
    app.Get("/", hi)

    // http://localhost:8080
    app.Listen(":8080")
}

type page struct {
    Title, Name string
}

func hi(ctx iris.Context) {
    //                      {{.Page.Title}} and {{Page.Name}}
    ctx.ViewData("Page", page{Title: "Hi Page", Name: "iris"})
    ctx.View("hi.html")
}
```

A real example can be found here: https://github.com/kataras/iris/tree/master/_examples/view/embedding-templates-into-app.

## Reload

Enable auto-reloading of templates on each request. Useful while developers are in dev mode
as they no neeed to restart their app on every template edit.

Example code:

```go
pugEngine := iris.Pug("./templates", ".jade")
pugEngine.Reload(true) // <--- set to true to re-build the templates on each request.
app.RegisterView(pugEngine)
```