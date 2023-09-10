---
title: Developing this blog in Go and HTMX
description: This article documents about how this blog came to be. From technical choices to deploying this blog.
---

## Motivation

I want to create, write and maintain a simple blog about personal and technical discoveries. I will self-host this blog on my Raspberry Pi cluster with kubernetes.

This already adds some constraints:

- Self-hostable.
- Lightweight: under 64Mo of RAM usage.
- Multiarch: ARM64 and AMD64.
- Low dependencies: for long term maintenance.
- Simplicity: simple articles, nothing special.
- Maintenable at long term.

These are the minimal constraints.

Because SSR permits partial update of the HTML documentation, I want SSR as a requirement to accelerate the rendering. Therefore, the requirements are:

- Writing and publishing articles in Markdown.
- SSR for routing.
- Index page.

## State of the art and Inspiration

There are a lot of existing solution to write a blog. I will only cite the solutions that are close to fill all requirements.

### Docusaurus

[Docusaurus](https://docusaurus.io) is a library on top of React which permits to write "content" very easily. It's light, fast, it has SSR and it is also extensible.

However, the drawbacks are:

- The dependencies: Too much dependencies. The latest [React is not even officially compatible](https://github.com/facebook/docusaurus/issues/7264#issuecomment-1257061642) with docusaurus. This prove basically that the product could involve a maintenance hell.
- React: This is my opinion, but, React is complex for nothing. The state management is hell, the component creation is conflicting (Class or function ? Man I don't care, just for one standard.). React may be "cool" for people who do a lot of web. development, to me it is just overkill.

### SvelteKit or SveltePress

[SvelteKit](https://kit.svelte.dev) is a blazingly fast JS framework for web development. It shines particularly at SSR, and it "fun" to write code with it. SvelteKit is "HTML"-first compared to its competitor [Next.js](https://nextjs.org).

[SveltePress](https://sveltepress.site) is a library on top of SvelteKit which simplifies the writing of "content" thanks to its supports for Markdown. A similar solution would be [Docusaurus](https://docusaurus.io), which is a library on top of React.

Thanks to [Vite](https://vitejs.dev), the resulting "bundle" is quite light and fast.

There isn't any drawbacks to say, and would have been my solution. And to be fair, I'm pretty sure this could have lead me to reach the end product a lot more faster.

So why I didn't choose this solution?

Because I don't want a framework to dictate how I should write my code. This is still overkill for a blog.

### HTML-only

Like the [motherfuckingwebsite](https://motherfuckingwebsite.com/), HTML-only would have been a pretty good solution. It's obviously SSR, It's readable, and simple. I mean, I can write HTML after all, so why not?

Simple: Markdown is better.

### Taking inspiration of my favorite blog structure: The Go dev blog.

[The Go Blog](https://go.dev/blog/) is a simple blog:

- There is an index, with titles, dates and authors.
- There are articles with content.

It uses the [Go HTTP server](https://cs.opensource.google/go/x/website/+/master:cmd/golangorg/server.go) to serve pages, and pages are in Markdown.

The only thing missing is SSR routing.

### Hugo

[Hugo](https://gohugo.io) is a framework to build static website focused on content. It uses Go under the hood. It is fast, light, etc... **SSR is missing obviously.** Also, why the [Hugo documentation is slow/weird](https://gohugo.io/documentation/)?

### HTMX and Go

[HTMX](https://htmx.org) is a JS library which abstracts the use of JS to manipulate the HTML documentation. It also allow the uses of Hypermedia as medium between the client and server. No need for JSON encode/decode. This is particularly great for SSR.

[Go](https://go.dev) is a programming language. It is compiled, has garbage collection, structures and is "OOP" compatible. Go stands out for its standard library, which allows easy concurrency and rapid development. Go's syntax is also fairly strict and explicit, so it's easy to read other developers' code. Basically, Go is dead simple and only allows one type of solution.

The reason I chose Go over C, Rust, C++, Java, ... is that static cross compilation is easy. Also, I write Go super fast and I don't have to fight with the language to choose a solution on "how I want to handle a string" (`String`, `char[]`, `std::string` ?!, give me one please!).

This is the solution I chose. Basically, I build my own HTTP server, similar to The Go Blog, and use HTMX to change pages quickly. I can take some of Hugo's Markdown rendering technologies, and bingo! I've got my blog!

Isn't just simple? (I'm gonna regret this.)

## Development

### Proof of Concept

Before even building a blog with HTMX, I must attest of HTMX capabilities. This is very important because, if a technology is too young, this means it is very sensible of breaking changes. When working actively on a project, breaking changes can be handled quite easily. However, this blog is a **side**-project, there are moment where I'm simply not maintaining this project, or worse, I haven't finished the project and simply throw this project to the trash.

Keep it simple stupid.

My PoC is to write a simple [authentication page](https://github.com/Darkness4/auth-htmx) with a counter behind it. It uses Go, HTMX and OAuth2. This is more complex than a blog, however, this PoC tests HTMX in depth.

My conclusion from this experiment is that:

- HTMX is pretty good for **only** rendering HTML. If I have a Android App, I would still need to use JSON or Protobuf.
- Error handling with HTMX is not that strightforward.
- HTMX is production-ready.

Which means it is pretty good for a writing blog.

### Architecture

#### Content directory

I'd like to follow the SvelteKit/Sveltepress directory structure. It's battle-tested, and sufficiently explicit.

Basically, there is a `pages` directory which accepts `.md` files and `.svelte` files. `.md` with [front matter](https://gohugo.io/content-management/front-matter/) are converted in html.

#### Templates and Components

With Go, templates can be named. With this, we can define components inside the `components` directory. These templates are passed to the templating engine.

#### Static directory

Like any HTTP server, there is a `static` directory where we can store unprocessed assets like images, icons, etc.

### Implementation

#### Initial Request and Server-Side rendering

The initial request happens when using the user uses the URL to seek a page. The initial request must load the HTMX library and application CSS. To do this, I've created a `base.html` template:

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

Clicking on `Go To Page 1` will make a HTMX request to the server `GET /page1` with the HTTP Header `Hx-Request: true`. The response of the server will replace the entire `body` element.

Thanks to that header, the server can identify if the client is doing an initial request or a SSR request:

```go
r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	var base string
	if r.Header.Get("Hx-Request") != "true" {
		// Initial Rendering
		base = "base.html"
	} else {
		// SSR
		base = "base.htmx"
	}

	// TODO: render template
})
```

The `base.htmx` is just simply the `body` and `head`:

```html
{{ define "base" }}
<head>
  {{ template "head" . }}
</head>

{{ template "body" . }} {{ end }}
```

#### Markdown rendering

Like Hugo, we will use the [Goldmark](https://github.com/yuin/goldmark) markdown rendering engine, coupled with chroma for syntax highlighting:

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

### Compile-time rendering

Since a blog is primarly static, I want to render the markdown and index page at compile-time.

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

As I said earlier, the HTTP server uses the `body` and `head` templates defined in the `gen/pages/blog/2023-09-09-hello-world/page.tmpl` for the initial request or SSR. This is how we can render the `body` and `head` templates with the `base` template:

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

```go
func buildPages() (index [][]Index, err error) {
	entries, err := os.ReadDir("gen/pages/blog")
	if err != nil {
		return index, err
	}

	// Sort the files in reverse order
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

With the implementation of these features, my blog has reached the minimum viable product. Missing from this article are discussions on:

- The Pager.
- Retrieving metadata (front matter) from an article in markdown.
- Escaping `{{` to avoid weird behavior with the templating engine.
- CSS, which is just PicoCSS.

The source code is available [here](https://github.com/Darkness4/blog).

Was it worth it? Hell yeah, I control everything in this blog, and most issues are indicated at compile time. With HTMX, PicoCSS and Goldmark, I don't write any HTML, CSS and JS, isn't that crazy?

Can this be a library like Hugo? F- no, I ain't ready for this.
