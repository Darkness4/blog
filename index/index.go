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
	"sort"
	"text/template"
	"time"

	"github.com/Darkness4/blog/utils/blog"
	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

// TODO: remove hardcoded value and put into a .env
const (
	title       = "Marc Nguyen's Blog"
	href        = "https://blog.mnguyen.fr"
	authorName  = "Marc Nguyen"
	authorEmail = "nguyen_marc@live.fr"
	description = "Marc Nguyen's blog is a personal and technical blog about documenting some processes, implementations, etc."
)

var created time.Time

func init() {
	created, _ = time.Parse("02 January 2006", "08 September 2023")
}

var (
	//go:embed templates/index.tmpl
	indexTmpl embed.FS
)

const elementPerPage = 50

type Index struct {
	Title         string
	Description   string
	PublishedDate int64
	Href          string
	EntryName     string
	Tags          []string
}

func buildPages() (index [][]Index, err error) {
	entries, err := os.ReadDir("gen/pages/blog")
	if err != nil {
		return index, err
	}

	// Sort the files in reverse order
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
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
		if !entry.IsDir() {
			continue
		}
		page := i / elementPerPage
		if page >= len(index) {
			index = append(index, make([]Index, 0, elementPerPage))
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
		tags := make([]string, 0, len(metaData["tags"].([]interface{})))
		for _, tag := range metaData["tags"].([]interface{}) {
			tags = append(tags, fmt.Sprintf("%v", tag))
		}
		index[page] = append(index[page], Index{
			EntryName:     entry.Name(),
			Title:         fmt.Sprintf("%v", metaData["title"]),
			Description:   fmt.Sprintf("%v", metaData["description"]),
			PublishedDate: date.Unix(),
			Href:          filepath.Join("/blog", entry.Name()),
			Tags:          tags,
		})
		i++
	}

	return index, nil
}

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
			Pages       [][]Index
			PageSize    int
			Title       string
			Href        string
			AuthorName  string
			AuthorEmail string
			Created     int64
			Updated     int64
			Description string
		}{
			Pages:       pages,
			PageSize:    len(pages),
			Updated:     time.Now().Unix(),
			Title:       title,
			Href:        href,
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			Created:     created.Unix(),
			Description: description,
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
