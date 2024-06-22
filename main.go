//go:generate go run -tags build build.go

package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"embed"

	"github.com/Darkness4/blog/db"
	"github.com/Darkness4/blog/gen/index"
	"github.com/Darkness4/blog/utils/color"
	"github.com/Darkness4/blog/utils/math"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	//go:embed gen components base.html base.htmx 404.tmpl
	html embed.FS
	//go:embed static
	static        embed.FS
	version       = "dev"
	listenAddress string
	publicURL     string
	dbDSN         string
)

var funcsMap = func() template.FuncMap {
	f := sprig.TxtFuncMap()
	f["computeColorByWord"] = color.ComputeByWord
	return f
}

func ReadUserIP(r *http.Request) string {
	IPAddress, _, _ := strings.Cut(r.Header.Get("X-Real-IP"), ",")
	if IPAddress == "" {
		IPAddress, _, _ = strings.Cut(r.Header.Get("X-Forwarded-For"), ",")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	addr := strings.ToLower(IPAddress)

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

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
		&cli.StringFlag{
			Name:        "db.dsn",
			Usage:       "The DSN for the database",
			Destination: &dbDSN,
			EnvVars:     []string{"DB_DSN"},
			Required:    true,
		},
	},
	Action: func(cCtx *cli.Context) error {
		ctx := cCtx.Context
		log.Level(zerolog.DebugLevel)

		// DB connection
		pool, err := pgxpool.New(ctx, dbDSN)
		if err != nil {
			return err
		}
		defer pool.Close()

		// Initial migration
		sqldb := stdlib.OpenDBFromPool(pool)
		if err := db.InitialMigration(sqldb); err != nil {
			return err
		}

		// Set up DB queries
		q := db.New(pool)

		// Router
		r := chi.NewRouter()
		r.Use(hlog.NewHandler(log.Logger))

		// Pages rendering
		var renderFn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			cleanPath := filepath.Clean(r.URL.Path)

			// Check if asset
			fpath := filepath.Join("gen/pages", cleanPath)
			if f, err := html.Open(fpath); err == nil {
				isPage := func() bool {
					defer f.Close()
					finfo, err := f.Stat()
					if err != nil {
						log.Err(err).Msg("failed to fetch fileinfo")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return false
					}

					if finfo.IsDir() {
						// It's a page, or a file not found
						return true
					}

					// Serve the file
					if _, err = io.Copy(w, f); err != nil {
						log.Err(err).Msg("failed to serve file")
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return false
				}()

				if !isPage {
					return
				}
			} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
				log.Err(err).Msg("failed to read file")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// It's a page
			templatePath := filepath.Clean(fmt.Sprintf("gen/pages/%s/page.tmpl", cleanPath))

			// Check if SSR
			var base string
			if r.Header.Get("Hx-Boosted") != "true" {
				// Initial Rendering
				base = "base.html"
			} else {
				// SSR
				base = "base.htmx"
			}

			// Set up the template
			t, err := template.New("base").
				Funcs(funcsMap()).
				ParseFS(html, base, templatePath, "components/*")
			var is404 bool
			if err != nil {
				if strings.Contains(err.Error(), "no files") {
					// Render 404
					w.WriteHeader(http.StatusNotFound)
					t, err = template.New("base").
						Funcs(funcsMap()).
						ParseFS(html, base, "404.tmpl", "components/*")
					if err != nil {
						panic(fmt.Sprintf("failed to parse 404.tmpl: %v", err))
					}
					is404 = true
				} else {
					log.Err(err).Msg("template error")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			pageS := r.URL.Query().Get("page")
			page, _ := strconv.Atoi(pageS)

			var pv db.PageView
			if !is404 {
				pv, err = q.FindPageViewsOrZero(ctx, strings.ToLower(cleanPath))
				if err != nil {
					log.Err(err).Msg("failed to fetch page views")
				}
			}

			if err := t.ExecuteTemplate(w, "base", struct {
				Pager struct {
					First   int
					Prev    int
					Current int
					Next    int
					Last    int
				}
				Index      []index.Index
				Path       string
				PublicURL  string
				PageViewsF string
				PageViews  int
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
				Index:      index.Pages[page],
				PageViewsF: math.FormatNumber(float64(pv.Views + 1)),
				PageViews:  int(pv.Views + 1),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !is404 {
				go func() {
					if err := q.CreateOrIncrementPageViewsOnUniqueIP(ctx, pool, strings.ToLower(cleanPath), ReadUserIP(r)); err != nil {
						log.Err(err).Msg("failed to increment page views")
					}
				}()
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
			fmt.Fprintf(w, `User-agent: *
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
	log.Logger = log.Logger.With().Caller().Logger()
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("app crashed")
	}
}
