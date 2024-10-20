---
title: Developing this blog in Go and HTMX
description: This article documents about how this blog came to be. From technical choices to deploying this blog.
tags: ['blog', 'go', 'htmx', 'raspberry-pi', 'kubernetes']
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

<hr>

## Motivation

I want to create, write and maintain a simple blog about personal and technical discoveries. I will self-host this blog on my Raspberry Pi cluster with Kubernetes.

This already adds some constraints:

- Self-hostable.
- Lightweight: under 64Mo of RAM usage.
- Multi-arch: ARM64 and AMD64.
- Low dependencies: for long term maintenance.
- Simplicity: simple articles, nothing special.
- Maintainable at long term.

These are the minimal constraints.

Because SSR permits partial update of the HTML documentation, I want SSR as a requirement to accelerate the rendering. Therefore, the requirements are:

- Writing and publishing articles in Markdown.
- SSR for routing.
- Index page.

## State of the art and Inspiration

There are a lot of existing solution to write a blog. I will only cite the solutions that are close to fill all requirements.

### Docusaurus

[Docusaurus](https://docusaurus.io) is a library on top of React that makes writing "content" very easy. It's lightweight, fast, has SSR and is also extensible.

However, the disadvantages are:

- Dependencies: Too many dependencies. The latest [React is not even officially compatible](https://github.com/facebook/docusaurus/issues/7264#issuecomment-1257061642) with Docusaurus. This basically could involve maintenance hell, since I cannot assure that each dependency (especially in the npm ecosystem) that any of them will be maintained.
- React: That's my opinion, but I think React is complex for nothing. State management is hellish, component creation is conflicting (class or function?!), and its performance is unsatisfactory. React may be "cool" for people who do a lot of web development, but for me, it's just an overkill. Alternatively, SvelteKit is already an excellent replacement for React by removing the VDOM and by having a compiler to check errors.

### SvelteKit or SveltePress

[SvelteKit](https://kit.svelte.dev) is an extremely fast JS framework for web development. It shines particularly at SSR, and is "fun" to write code with. SvelteKit is more "HTML"-first than its competitor [Next.js](https://nextjs.org), React, etc.

[SveltePress](https://sveltepress.site) is a library on top of SvelteKit which simplifies the writing of "content" thanks to its supports for Markdown, an equivalent of Docusaurus.

Thanks to [Vite](https://vitejs.dev), the resulting "bundle" is quite light and fast.

There's no downside to saying, and that would have been my solution. And to be honest, I'm pretty sure it would have got me to the final product a lot quicker.

Why didn't I choose this solution?

Because I don't want a framework telling me how to write my code. It's still too much for a blog.

### HTML-only

As with [motherfuckingwebsite](https://motherfuckingwebsite.com/), HTML alone would have been a good solution. It's obviously SSR, it's readable, and it's simple. I mean, I can write HTML after all, so why not?

It's simple: Markdown is better.

### Taking inspiration of my favorite blog structure: The Go dev blog.

[The Go Blog](https://go.dev/blog/) is a simple blog:

- There is an index, with titles, dates and authors.
- There are articles with content.

It uses the [Go HTTP server](https://cs.opensource.google/go/x/website/+/master:cmd/golangorg/server.go) to serve pages, and pages are in Markdown.

The only thing missing is SSR routing.

### Hugo

[Hugo](https://gohugo.io) is a framework for building static, content-centric websites. It uses Go under the hood. It's fast, lightweight, etc... **Also, why is [Hugo documentation slow/weird](https://gohugo.io/documentation/)**?

Anyway, it's a pass for me.

### HTMX and Go

[HTMX](https://htmx.org) is a JS library that abstracts the use of JS to manipulate HTML documentation. It also allows the use of hypermedia as a medium between client and server. There's no need to encode/decode JSON. This is particularly interesting for SSR where latency is important.

[Go](https://go.dev) is a programming language. It is compiled, has a garbage collection system, structures and is OOP-compatible. Go stands out for its standard library, which allows easy concurrency and rapid development. Go's syntax is also fairly strict and explicit, so it's easy to read other developers' code. In fact, Go is so simple that most of the time it allows only one type of solution, thus achieving the [Zen of Python](https://peps.python.org/pep-0020/) better than Python itself.

The reason I chose Go over C, Rust, C++, Java, ... is that static cross compilation is easy. Also, I write Go superfast and I don't have to fight with the language to choose a solution on "how I want to handle a string" (`String`, `char[]`, `std::string` ?!, give me one please!).

This is the solution I've chosen. Basically, the idea is as follows:

- I build my own HTTP server, similar to The Go Blog
- I use HTMX to change pages quickly.
- I can use some of Hugo's Markdown rendering technologies.

And bingo! I've got my blog!

Simple, isn't it? (I'm going to regret this.)

## Development

### Proof of Concept

Before I even start building a blog with HTMX, I have to certify HTMX's capabilities. This is very important because, if a technology is too young, it means it's very sensitive to breaking changes. When you're actively working on a project, disruptions can be managed quite easily. However, this blog is a **side project**, there are times when I simply don't maintain this project, or worse, I haven't completed the project and simply throw it in the garbage can.

Keep it simple stupid.

My PoC is to write a simple [authentication page](https://github.com/Darkness4/auth-htmx) with a counter behind it. It uses Go, HTMX and OAuth2. It's more complex than a blog, but this PoC tests HTMX in depth.

Testing with OAuth2 is a good PoC because it tests:

- The UX of login flow (login, logout).
- A small state management ("is logged").
- The authorization ("count only if is logged").

The conclusion I draw from this experiment is as follows:

- HTMX is good enough for **only** rendering HTML. If I had an Android application, I'd still need to use JSON or Protobuf.
- Error handling with HTMX is not that simple. Some say it's just different, which I can agree.
- HTMX is ready for production.

Which means it's great for writing a blog!

### Architecture

#### Content directory

I followed the SvelteKit/Sveltepress directory structure. It's battle-tested, and sufficiently explicit.

Basically, there is a `pages` directory which accepts `.md` files and `.svelte` files. `.md` with [front matter](https://gohugo.io/content-management/front-matter/) are converted in html using a Markdown rendering engine.

#### Templates and Components

With Go, templates can be named. With this, we can define components inside the `components` directory. These templates are passed to the template engine.

There is also the `base.html` and `base.htmx` templates for the initial request and SSR request with HTMX.

The final template is `markdown.tmpl`, which is used with the Markdown renderer to arrange the page. It essentially surrounds the output of the Markdown renderer.

#### Static directory

Like any HTTP server, there is a `static` directory where we can store unprocessed assets like images, icons, etc. This is served like any Go file server.

### Implementation

#### Initial Request and Server-Side rendering

The initial request happens when using the user uses the URL to seek a page. The initial request must load the HTMX library and application CSS. To do this, I've created a `base.html` template.

**base.html**

```html
{{ define "base" }}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta hx-preserve="true" charset="UTF-8" />
    <meta
      hx-preserve="true"
      name="viewport"
      content="width=device-width, initial-scale=1"
    />
    <script hx-preserve="true" src="https://unpkg.com/htmx.org@1.9.5"></script>
    <script
      hx-preserve="true"
      src="https://unpkg.com/htmx.org@1.9.5/dist/ext/head-support.js"
    ></script>
    <link
      hx-preserve="true"
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/@picocss/pico@next/css/pico.classless.min.css"
    />
    <link hx-preserve="true" rel="stylesheet" href="/static/app.css" />
    {{ template "head" . }}
  </head>
  <body hx-ext="head-support" hx-boost="true">
    {{ template "body" . }}
  </body>
</html>
{{ end }}
```

This is very similar to SvelteKit's [`app.html`](https://github.com/sveltejs/kit/blob/92eb1aacb2355f36dd0d0e88048b8a9bd8219ea3/packages/create-svelte/templates/default/src/app.html).

Pages in the `pages` directory define the `head` and `body` template. Therefore, when reaching any URL, the template will be rendered base on the path

When the user has loaded the initial page, clicking on any `a` element will make a HTMX request to server thanks to [`hx-boost`](https://htmx.org/attributes/hx-boost/).

The `hx-boost` attribute will replace the `body` when the user clicks on a `a` element. I've also added the `head-support` for HTMX so that we can also replace `head` when doing SSR. This is used to dynamically change the `title` and CSS.

Example:

```html
<div hx-boost="true">
  <a href="/page1">Go To Page 1</a>
  <a href="/page2">Go To Page 2</a>
</div>
```

Clicking on `Go To Page 1` will make a HTMX request to the server `GET /page1` with the HTTP Header `Hx-Boosted: true`. The response of the server will replace the entire `body` element.

Thanks to the header `Hx-Boosted: true`, the server can identify if the client is doing an initial request or an SSR request:

**main.go**

```go
r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	var base string
	if r.Header.Get("Hx-Boosted") != "true" {
		// Initial Rendering
		base = "base.html"
	} else {
		// SSR
		base = "base.htmx"
	}

	// TODO: render template
})
```

The `base.htmx` is just simply the `body` and `head`.

**base.htmx**

```html
{{ define "base" }}
<head>
  {{ template "head" . }}
</head>

{{ template "body" . }} {{ end }}
```

#### Markdown rendering

Like Hugo, we will use the [Goldmark](https://github.com/yuin/goldmark) markdown rendering engine, coupled with chroma for syntax highlighting:

**Example:**

```go
var cssBuffer strings.Builder
markdown := goldmark.New(
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("onedark"),
			highlighting.WithCSSWriter(&cssBuffer),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.WithClasses(true),
			),
		),
		extension.GFM,
		meta.Meta,
		&anchor.Extender{},
	),
)

var out strings.Builder
ctx := parser.NewContext()
err := markdown.Convert(content, &out, parser.WithContext(ctx))
```

The resulting string can be passed as the body:

```go
if err := t.Execute(w, struct {
	Style string
	Body  string
}{
	Style: cssBuffer.String(),
	Body:  out.String(),
}); err != nil {
	log.Fatal().Err(err).Msg("generate file from template failure")
}
```

#### Compile-time rendering

Since a blog is primarily static, I want to render the markdown and index page at compile-time.

Go doesn't have any [`comptime`](https://ziglang.org/documentation/master/#comptime) like Zig, but we can at least use `go generate`.

The idea is to create a Go "script" which can be executed by `go generate`.

**build.go**:

```go
//go:build build

package main

func main() {
	// TODO: compile time stuff
	processPages()
	index.Generate()
}
```

**main.go**:

```go
//go:generate go run -tags build build.go

package main

func main() {
	// TODO: runtime time stuff
}
```

The `processPages` function executes the markdown rendering and outputs the resulting files in the `gen/` directory.

Example: `pages/blog/2023-09-09-hello-world/page.md` ⟶ `gen/pages/blog/2023-09-09-hello-world/page.tmpl`.

As I said earlier, the HTTP server uses the `body` and `head` templates defined in the `gen/pages/blog/2023-09-09-hello-world/page.tmpl` for the initial request or SSR. This is how we can render the `body` and `head` templates with the `base` template.

**main.go**

```go
//go:embed gen components base.html base.htmx
var html embed.FS

// ... in the main function
r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	// ...
	// base = "base.html" or "base.htmx"
	path := filepath.Clean(r.URL.Path)
	// TODO: check if it's page or an asset
	path = filepath.Clean(fmt.Sprintf("gen/pages/%s/page.tmpl", path))

	t, err := template.New("base").
		Funcs(sprig.TxtFuncMap()).
		ParseFS(html, base, path, "components/*")
	if err != nil {
		if strings.Contains(err.Error(), "no files") {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			log.Err(err).Msg("template error")
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	if err := t.ExecuteTemplate(w, "base", struct {}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
})
```

#### Index page and pagination

This page is also generated at compile-time.

Therefore, at compile-time, we need a function that fetches the list of blog pages, and group them.

**index.go**

```go
func buildPages() (index [][]Index, err error) {
	// Read the blog pages
	entries, err := os.ReadDir("gen/pages/blog")
	if err != nil {
		return index, err
	}

	// Sort the pages in reverse order
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() > entries[i].Name()
	})

	// Markdown Parser
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)

	index = make([][]Index, 0, len(entries)/elementPerPage+1)
	i := 0
	for _, entry := range entries {
		// Should be a dir
		if !entry.IsDir() {
			continue
		}
		page := i / elementPerPage
		if page >= len(index) {
			index = append(index, make([]Index, 0, elementPerPage))
		}
		// Should contains a page.md
		f, err := os.Open(filepath.Join("pages/blog", entry.Name(), "page.md"))
		if err != nil {
			log.Debug().
				Err(err).
				Str("entry", entry.Name()).
				Msg("ignored for index, failed to open page.md")
			continue
		}
		finfo, err := f.Stat()
		if err != nil {
			log.Debug().
				Err(err).
				Str("entry", entry.Name()).
				Msg("ignored for index, failed to stat page.md")
			continue
		}
		if finfo.IsDir() {
			continue
		}

		// Fetch metadata
		b, err := io.ReadAll(f)
		if err != nil {
			log.Fatal().Err(err).Msg("read file failure")
		}
		document := markdown.Parser().Parse(text.NewReader(b))
		metaData := document.OwnerDocument().Meta()

		index[page] = append(index[page], Index{
			EntryName:     entry.Name(),
			Title:         fmt.Sprintf("%v", metaData["title"]),
			Description:   fmt.Sprintf("%v", metaData["description"]),
			PublishedDate: date.Format("Monday 02 January 2006"),
			Href:          filepath.Join("/blog", entry.Name()),
		})
		i++
	}

	return index, nil
}
```

Then, we generate a `.go` that contains the index:

**index.tmpl**

```go
{{define "index"}}
package index

type Index struct {
	Title         string
	Description   string
	Href          string
	EntryName     string
}

const PageSize = {{ .PageSize }}

var Pages = [][]Index{
	{{- range $page := .Pages}}
	{
		{{- range $i, $value := $page}}
		{
			EntryName: {{ $value.EntryName | quote }},
			Title: {{ $value.Title | quote }},
			Description: {{ $value.Description | quote }},
			Href: {{ $value.Href | quote }},
		},
		{{- end}}
	},
	{{- end}}
}
{{- end}}


```

**index.go** (used to render the template)

```go
func Generate() {
	pages, err := buildPages()
	if err != nil {
		log.Fatal().Err(err).Msg("index failure")
	}

	out := "gen/index/index.go"
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		log.Fatal().Err(err).Msg("mkdir failure")
	}

	func() {
		f, err := os.Create(out)
		if err != nil {
			log.Fatal().Err(err).Msg("generate file from template failure")
		}
		defer f.Close()

		var buf bytes.Buffer
		t := template.Must(
			template.New("index").
				Funcs(sprig.TxtFuncMap()).
				ParseFS(indexTmpl, "templates/index.tmpl"),
		)
		if err := t.ExecuteTemplate(&buf, "index", struct {
			Pages    [][]Index
			PageSize int
		}{
			Pages:    pages,
			PageSize: len(pages),
		}); err != nil {
			log.Fatal().Err(err).Msg("template failure")
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println(buf.String())
			log.Fatal().Err(err).Msg("format code from template failure")
		}

		if _, err = f.Write(formatted); err != nil {
			log.Fatal().Err(err).Msg("write failure")
		}
	}()
}

```

## Conclusion

To summary, when I write a page of the blog, I create a `page.md` file in the `pages/blog/<date-url>/` directory like SveltePress. Then I compile the pages into html by running `go generate` which `go run build.go`. The `build.go` script generates the html pages from markdown, and also creates the `gen/index/index.go` file which contains the necessary data for the index page.

After compiling the pages, the server can run and accepts two type of request: the initial request and the SSR request with HTMX. Based on the request, we can create a router which executes SSR and allow quick rendering without the need to redownload everything (like the HTMX library, CSS, and other "app"-time assets...).

With the implementation of these features, my blog has reached the minimum viable product. Missing from this article are discussions on:

- The Pager.
- Escaping `{{` to avoid weird behavior with the templating engine.
- CSS, which is just PicoCSS.

_Was it worth it?_ Hell yeah, I control everything in this blog, and most problems are indicated at compile time. With HTMX, PicoCSS and Goldmark, I'm no longer writing HTML, CSS and JS to write an article, how crazy is that?

_How difficult was it?_ It's not really "hard", but there are certainly more steps than bootstrapping SveltePress or Docusaurus. It's also more prone to bugs than an established product like SveltePress or Docusaurus.

_Can it be a library like Hugo?_ No, I'm not ready for that.

Anyway, this is like installing Gentoo: it's for learning, having full-control and being ultra-optimized.

## References

- [Blog's Source Code](https://github.com/Darkness4/blog)
- [HTMX](https://htmx.org)
- [Goldmark](https://github.com/yuin/goldmark)
- [PicoCSS](https://picocss.com)
