//go:build build

package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
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

func main() {
	// Markdown engine
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("onedark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			meta.Meta,
		),
	)

	for file := range files() {
		content, err := md.ReadFile(file)
		if err != nil {
			log.Fatal().Err(err).Msg("read file failure")
		}
		ext := filepath.Ext(file)
		file = filepath.Join("gen", strings.TrimSuffix(file, ext)+".tmpl")

		if err := os.MkdirAll(filepath.Dir(file), 0o777); err != nil {
			log.Fatal().Err(err).Msg("mkdir failure")
		}

		func() {
			w, err := os.Create(file)
			if err != nil {
				log.Fatal().Err(err).Msg("create file failure")
			}
			defer w.Close()

			if ext == ".md" {
				var sb strings.Builder

				ctx := parser.NewContext()
				if err := markdown.Convert(content, &sb, parser.WithContext(ctx)); err != nil {
					log.Fatal().Err(err).Msg("write file failure")
				}
				metaData := meta.Get(ctx)

				t := template.Must(template.ParseFS(mdTmpl, "templates/markdown.tmpl"))
				t.Execute(w, struct {
					Title string
					Body  string
				}{
					Title: fmt.Sprintf("%v", metaData["title"]),
					Body:  sb.String(),
				})
			} else {
				if _, err := w.Write(content); err != nil {
					log.Fatal().Err(err).Msg("write file failure")
				}
			}
		}()
	}
}
