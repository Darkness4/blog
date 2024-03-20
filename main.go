//go:generate go run -tags build build.go

package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"embed"

	"github.com/Darkness4/blog/gen/index"
	"github.com/Darkness4/blog/utils/math"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed gen components base.html base.htmx 404.tmpl
	html embed.FS
	//go:embed static
	static        embed.FS
	version       = "dev"
	listenAddress string
	publicURL     string
)

var app = &cli.App{
	Name:    "blog",
	Version: version,
	Usage:   "A blog in HTMX.",
	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "listen.address",
			Usage:       "The address to listen on",
			Value:       ":3000",
			Destination: &listenAddress,
			EnvVars:     []string{"LISTEN_ADDRESS"},
		},
		&cli.StringFlag{
			Name:        "public.url",
			Usage:       "The public URL",
			Value:       "https://blog.mnguyen.fr",
			Destination: &publicURL,
			EnvVars:     []string{"PUBLIC_URL"},
		},
	},
	Action: func(_ *cli.Context) error {
		log.Level(zerolog.DebugLevel)

		// Router
		r := chi.NewRouter()
		r.Use(hlog.NewHandler(log.Logger))

		// Pages rendering
		var renderFn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Clean(r.URL.Path)

			// Check if asset
			if f, err := html.Open(filepath.Join("gen/pages", r.URL.Path)); err == nil {
				finfo, err := f.Stat()
				if err != nil {
					log.Err(err).Msg("failed to fetch fileinfo")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if !finfo.IsDir() {
					buf := make([]byte, 512)
					for {
						_, err := f.Read(buf)
						if err == io.EOF {
							return
						} else if err != nil {
							log.Err(err).Msg("failed to read body")
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						if _, err := w.Write(buf); err != nil {
							log.Err(err).Msg("failed to write body")
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
				}
				_ = f.Close()
			} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
				log.Err(err).Msg("failed to read file")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// It's a page
			path = filepath.Clean(fmt.Sprintf("gen/pages/%s/page.tmpl", path))

			// Check if SSR
			var base string
			if r.Header.Get("Hx-Boosted") != "true" {
				// Initial Rendering
				base = "base.html"
			} else {
				// SSR
				base = "base.htmx"
			}
			t, err := template.New("base").
				Funcs(sprig.TxtFuncMap()).
				ParseFS(html, base, path, "components/*")
			if err != nil {
				if strings.Contains(err.Error(), "no files") {
					w.WriteHeader(http.StatusNotFound)
					t, err = template.New("base").
						Funcs(sprig.TxtFuncMap()).
						ParseFS(html, base, "404.tmpl", "components/*")
					if err != nil {
						panic(fmt.Sprintf("failed to parse 404.tmpl: %v", err))
					}
				} else {
					log.Err(err).Msg("template error")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			pageS := r.URL.Query().Get("page")
			page, _ := strconv.Atoi(pageS)
			if err := t.ExecuteTemplate(w, "base", struct {
				Pager struct {
					First   int
					Prev    int
					Current int
					Next    int
					Last    int
				}
				Index     []index.Index
				Path      string
				PublicURL string
			}{
				PublicURL: publicURL,
				Path:      r.URL.Path,
				Pager: struct {
					First   int
					Prev    int
					Current int
					Next    int
					Last    int
				}{
					First:   0,
					Prev:    math.MaxI(0, page-1),
					Current: page,
					Next:    math.MinI(index.PageSize-1, page+1),
					Last:    index.PageSize - 1,
				},
				Index: index.Pages[page],
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		r.Get("/rss", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/rss+xml")
			if err := index.Feed.WriteRss(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/atom", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/atom+xml")
			if err := index.Feed.WriteAtom(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/json", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := index.Feed.WriteJSON(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		r.Get("/sitemap.xml", func(w http.ResponseWriter, _ *http.Request) {
			b, err := index.ToSiteMap(index.Pages)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/xml")
			header := `<?xml version="1.0" encoding="UTF-8"?>`
			fmt.Fprintf(w, "%s\n%s", header, b)
		})
		r.Get("/robots.txt", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, `User-agent: *
Disallow:

Sitemap: %s/sitemap.xml
Sitemap: %s/rss
Sitemap: %s/atom
`, publicURL, publicURL, publicURL)
		})
		r.Get("/*", renderFn)
		r.Handle("/static/*", http.FileServer(http.FS(static)))

		log.Info().Str("listenAddress", listenAddress).Msg("listening")
		return http.ListenAndServe(listenAddress, r)
	},
}

func main() {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("app crashed")
	}
}
