//go:build build

package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Darkness4/blog/d2"
	"github.com/Darkness4/blog/images"
	"github.com/Darkness4/blog/index"
	"github.com/Darkness4/blog/markdown"
	"github.com/Darkness4/blog/utils/blog"
	"github.com/Darkness4/blog/utils/ptr"
	"github.com/Darkness4/blog/utils/unique"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/rs/zerolog/log"
	admonitions "github.com/stefanfritsch/goldmark-admonitions"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/toc"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
)

var (
	//go:embed pages/*
	md embed.FS

	//go:embed templates/markdown.tmpl
	mdTmpl embed.FS
)

func processDirectory(fs embed.FS, dirPath string, filePaths chan<- string) error {
	out, err := fs.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range out {
		filePath := filepath.Join(dirPath, file.Name())

		if file.IsDir() {
			// If it's a directory, recursively process it
			if err := processDirectory(fs, filePath, filePaths); err != nil {
				return err
			}
		} else {
			filePaths <- filePath
		}
	}

	return nil
}

func files() <-chan string {
	f := make(chan string, 1)
	go func() {
		defer close(f)
		processDirectory(md, "pages", f)
	}()
	return f
}

func filterBlogPages(input <-chan string) (blogPages <-chan string, rest <-chan string) {
	o := make(chan string, 100)
	r := make(chan string, 100)
	go func() {
		defer close(o)
		defer close(r)
		for file := range input {
			if strings.HasPrefix(file, "pages/blog") && strings.HasSuffix(file, "page.md") {
				o <- file
			} else {
				r <- file
			}
		}
	}()
	return o, r
}

func triple(input <-chan string) <-chan struct {
	prev string
	curr string
	next string
} {
	output := make(chan struct {
		prev string
		curr string
		next string
	}, 100)
	go func() {
		defer close(output)
		var prev, curr string
		for next := range input {
			if curr != "" {
				output <- struct {
					prev string
					curr string
					next string
				}{prev: prev, curr: curr, next: next}
			}
			prev, curr = curr, next
		}

		if curr != "" {
			output <- struct {
				prev string
				curr string
				next string
			}{prev: prev, curr: curr, next: ""}
		}
	}()
	return output
}

func processPages() {
	_ = os.RemoveAll("gen")

	// Markdown engine
	cssBuffer := unique.NewLineWriter()
	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
			renderer.WithNodeRenderers(util.Prioritized(markdown.NewRenderer(), 1)),
		),
		images.NewReplacer(func(link string) string {
			if filepath.IsAbs(link) || strings.HasPrefix(strings.ToLower(link), "http") {
				return link
			}
			return filepath.Join("\\{\\{ $.Path }}", link)
		}),
		goldmark.WithExtensions(
			mathjax.MathJax,
			&d2.Extender{
				RenderOptions: d2.RenderOptions{
					ThemeID: &d2themescatalog.DarkMauve.ID,
					Scale:   ptr.Ref(0.9),
					Pad:     ptr.Ref(int64(25)),
					Center:  ptr.Ref(true),
				},
			},
			highlighting.NewHighlighting(
				highlighting.WithStyle("onedark"),
				highlighting.WithCSSWriter(cssBuffer),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.WithClasses(true),
				),
			),
			extension.GFM,
			meta.Meta,
			&anchor.Extender{
				Texter: anchor.Text("#"),
				Attributer: anchor.Attributes{
					"class": "anchor",
				},
			},
			&admonitions.Extender{},
		),
	)

	files := files()
	blogPages, files := filterBlogPages(files)

	// Process blogPages
	for file := range triple(blogPages) {
		curr := file.curr
		content, err := md.ReadFile(curr)
		if err != nil {
			log.Fatal().Err(err).Msg("read file failure")
		}
		ext := filepath.Ext(curr)
		curr = filepath.Join("gen", strings.TrimSuffix(curr, ext))

		if err := os.MkdirAll(filepath.Dir(curr), 0o755); err != nil {
			log.Fatal().Err(err).Msg("mkdir failure")
		}

		func() {
			w, err := os.Create(curr + ".tmpl")
			if err != nil {
				log.Fatal().Err(err).Msg("create file failure")
			}
			defer w.Close()

			var sb strings.Builder

			ctx := parser.NewContext()
			doc := markdown.Parser().Parse(text.NewReader(content), parser.WithContext(ctx))
			if err := markdown.Renderer().Render(&sb, content, doc); err != nil {
				log.Fatal().Err(err).Msg("write file failure")
			}
			tree, err := toc.Inspect(doc, content, toc.Compact(true))
			if err != nil {
				log.Fatal().Err(err).Msg("toc failure")
			}
			var tocSB strings.Builder
			list := toc.RenderList(tree)
			if list != nil {
				if err := markdown.Renderer().Render(&tocSB, content, list); err != nil {
					log.Fatal().Err(err).Msg("toc render failure")
				}
			}

			// Replace {{ with &#123;&#123;
			out := strings.ReplaceAll(sb.String(), "{{", "&#123;&#123;")
			out = strings.ReplaceAll(out, "\\{\\{", "{{")
			metaData := meta.Get(ctx)
			date, err := blog.ExtractDate(filepath.Base(filepath.Dir(curr)))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse date failure")
			}

			// Compile time variable
			var bodySB strings.Builder
			tBody := template.Must(template.New("body").Parse(out))
			if err := tBody.Execute(&bodySB, struct {
				TOC  string
				Path string
			}{
				TOC:  tocSB.String(),
				Path: "{{ $.Path }}", // Pass variable to runtime
			}); err != nil {
				log.Fatal().Err(err).Msg("body template failure")
			}

			t := template.Must(template.ParseFS(mdTmpl, "templates/markdown.tmpl"))
			if err := t.Execute(w, struct {
				Title         string
				Description   string
				Style         string
				Body          string
				PublishedDate string
				TOC           string

				Prev string
				Next string
			}{
				Title:         fmt.Sprintf("%v", metaData["title"]),
				Description:   fmt.Sprintf("%v", metaData["description"]),
				Style:         cssBuffer.String(),
				Body:          bodySB.String(),
				TOC:           tocSB.String(),
				PublishedDate: date.Format("Monday 02 January 2006"),

				Prev: strings.TrimSuffix(strings.TrimPrefix(file.prev, "pages"), "/page.md"),
				Next: strings.TrimSuffix(strings.TrimPrefix(file.next, "pages"), "/page.md"),
			}); err != nil {
				log.Fatal().Err(err).Msg("generate file from template failure")
			}
			cssBuffer.Reset()
		}()
	}

	for file := range files {
		content, err := md.ReadFile(file)
		if err != nil {
			log.Fatal().Err(err).Msg("read file failure")
		}
		ext := filepath.Ext(file)
		file = filepath.Join("gen", strings.TrimSuffix(file, ext))

		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			log.Fatal().Err(err).Msg("mkdir failure")
		}

		func() {
			if ext == ".md" {
				w, err := os.Create(file + ".tmpl")
				if err != nil {
					log.Fatal().Err(err).Msg("create file failure")
				}
				defer w.Close()

				var sb strings.Builder

				ctx := parser.NewContext()
				if err := markdown.Convert(content, &sb, parser.WithContext(ctx)); err != nil {
					log.Fatal().Err(err).Msg("write file failure")
				}
				metaData := meta.Get(ctx)
				out := strings.ReplaceAll(sb.String(), "{{", "&#123;&#123;")

				t := template.Must(template.ParseFS(mdTmpl, "templates/markdown.tmpl"))
				if err := t.Execute(w, struct {
					Title       string
					Description string
					Style       string
					Body        string
				}{
					Title:       fmt.Sprintf("%v", metaData["title"]),
					Description: fmt.Sprintf("%v", metaData["description"]),
					Style:       cssBuffer.String(),
					Body:        out,
				}); err != nil {
					log.Fatal().Err(err).Msg("generate file from template failure")
				}
				cssBuffer.Reset()
			} else {
				w, err := os.Create(file + ext)
				if err != nil {
					log.Fatal().Err(err).Msg("create file failure")
				}
				defer w.Close()

				if _, err := w.Write(content); err != nil {
					log.Fatal().Err(err).Msg("write file failure")
				}
			}
		}()
	}
}

func main() {
	processPages()
	index.Generate()
}
