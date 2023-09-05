//go:build build

package index

import (
	"bytes"
	"embed"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

var (
	//go:embed templates/index.tmpl
	indexTmpl embed.FS
)

type Index struct {
}

func buildIndex() (index map[string]Index, err error) {
	entries, err := os.ReadDir("gen/pages/blog")
	if err != nil {
		return index, err
	}

	// Filters non-page
	index = make(map[string]Index)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join("gen/pages/blog", entry.Name(), "page.tmpl"))
		if err != nil {
			log.Debug().
				Err(err).
				Str("entry", entry.Name()).
				Msg("ignored for index, failed to open page.tmpl")
			continue
		}
		finfo, err := f.Stat()
		if err != nil {
			log.Debug().
				Err(err).
				Str("entry", entry.Name()).
				Msg("ignored for index, failed to stat page.tmpl")
			continue
		}
		if finfo.IsDir() {
			continue
		}
		index[entry.Name()] = Index{}
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
		t := template.Must(template.ParseFS(indexTmpl, "templates/index.tmpl"))
		if err := t.Execute(&buf, struct {
			Module string
			Index  map[string]Index
		}{
			Module: bi.Deps[0].Path,
			Index:  index,
		}); err != nil {
			log.Fatal().Err(err).Msg("template failure")
		}

		formated, err := format.Source(buf.Bytes())
		if err != nil {
			log.Fatal().Err(err).Msg("format code from template failure")
		}

		f.Write(formated)
	}()
}
