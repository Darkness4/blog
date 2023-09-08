//go:build build

package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Darkness4/blog/index"
	"github.com/Darkness4/blog/utils/blog"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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
	var cssBuffer strings.Builder
	markdown := goldmark.New(
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
			if err := markdown.Convert(content, &sb, parser.WithContext(ctx)); err != nil {
				log.Fatal().Err(err).Msg("write file failure")
			}
			metaData := meta.Get(ctx)
			date, err := blog.ExtractDate(filepath.Base(filepath.Dir(curr)))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse date failure")
			}

			t := template.Must(template.ParseFS(mdTmpl, "templates/markdown.tmpl"))
			if err := t.Execute(w, struct {
				Title         string
				Style         string
				Body          string
				PublishedDate string

				Prev string
				Next string
			}{
				Title:         fmt.Sprintf("%v", metaData["title"]),
				Style:         cssBuffer.String(),
				Body:          sb.String(),
				PublishedDate: date.Format("Monday 02 January 2006"),

				Prev: strings.TrimSuffix(strings.TrimPrefix(file.next, "pages"), "/page.md"),
				Next: strings.TrimSuffix(strings.TrimPrefix(file.prev, "pages"), "/page.md"),
			}); err != nil {
				log.Fatal().Err(err).Msg("generate file from template failure")
			}
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

				t := template.Must(template.ParseFS(mdTmpl, "templates/markdown.tmpl"))
				if err := t.Execute(w, struct {
					Title string
					Style string
					Body  string
				}{
					Title: fmt.Sprintf("%v", metaData["title"]),
					Style: cssBuffer.String(),
					Body:  sb.String(),
				}); err != nil {
					log.Fatal().Err(err).Msg("generate file from template failure")
				}
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
