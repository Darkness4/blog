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

	"github.com/Darkness4/blog/utils/blog"
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

const elementPerPage = 50

type Index struct {
	Title         string
	Description   string
	PublishedDate string
	Href          string
}

func buildPages() (index []map[string]Index, err error) {
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
	index = make([]map[string]Index, len(entries)/elementPerPage+1)
	for i, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		page := i / elementPerPage
		if index[page] == nil {
			index[page] = make(map[string]Index)
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
		date, err := blog.ExtractDate(entry.Name())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read date")
		}

		index[page][entry.Name()] = Index{
			Title:         fmt.Sprintf("%v", metaData["title"]),
			Description:   fmt.Sprintf("%v", metaData["description"]),
			PublishedDate: date.Format("Monday 02 January 2006"),
			Href:          filepath.Join("/blog", entry.Name()),
		}
	}

	return index, nil
}

func Generate() {
	pages, err := buildPages()
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
			Module   string
			Pages    []map[string]Index
			PageSize int
		}{
			Module:   bi.Deps[0].Path,
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
