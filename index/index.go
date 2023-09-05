//go:build build

package index

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

var (
	//go:embed templates/index.tmpl
	indexTmpl embed.FS
)

type Index struct {
	Title       string
	Description string
	Href        string
}

func buildIndex() (index map[string]Index, err error) {
	entries, err := os.ReadDir("gen/pages/blog")
	if err != nil {
		return index, err
	}

	// Markdown Parser
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)

	// Filters non-page
	index = make(map[string]Index)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
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
		b, err := io.ReadAll(f)
		if err != nil {
			log.Fatal().Err(err).Msg("read file failure")
		}
		document := markdown.Parser().Parse(text.NewReader(b))
		metaData := document.OwnerDocument().Meta()
		index[entry.Name()] = Index{
			Title:       fmt.Sprintf("%v", metaData["title"]),
			Description: fmt.Sprintf("%v", metaData["description"]),
			Href:        filepath.Join("blog", entry.Name()),
		}
	}

	return index, nil
}

func Generate() {
	index, err := buildIndex()
	if err != nil {
		log.Fatal().Err(err).Msg("index failure")
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info")
		return
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
			Module string
			Index  map[string]Index
		}{
			Module: bi.Deps[0].Path,
			Index:  index,
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
