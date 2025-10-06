package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Darkness4/blog/db"
	"github.com/Darkness4/blog/web"
	"github.com/Darkness4/blog/web/gen/index"
	"github.com/Darkness4/blog/web/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	version       = "dev"
	listenAddress string
	publicURL     string
	dbDSN         string
	csp           string
)

var app = &cli.Command{
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
			Sources:     cli.EnvVars("LISTEN_ADDRESS"),
		},
		&cli.StringFlag{
			Name:        "public.url",
			Usage:       "The public URL",
			Value:       "https://blog.mnguyen.fr",
			Destination: &publicURL,
			Sources:     cli.EnvVars("PUBLIC_URL"),
		},
		&cli.StringFlag{
			Name:        "db.dsn",
			Usage:       "The DSN for the database",
			Destination: &dbDSN,
			Sources:     cli.EnvVars("DB_DSN"),
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "csp",
			Usage:       "The Content Security Policy",
			Destination: &csp,
			Sources:     cli.EnvVars("CSP"),
			Value: `default-src 'self';
base-uri 'self';
form-action 'self';
frame-ancestors 'none';
script-src 'self' 'unsafe-inline' https://unpkg.com/ https://cloud.umami.is/;
style-src 'self' 'unsafe-inline' https://unpkg.com/ https://fonts.googleapis.com/;
connect-src 'self' https://cloud.umami.is/;
media-src 'self' https://www.youtube.com/ https://www.youtube-nocookie.com/;
frame-src https://www.youtube.com/ https://www.youtube-nocookie.com/;
font-src 'self' https://fonts.gstatic.com/;
img-src 'self' data: *;
object-src 'none';`,
		},
	},
	Action: func(ctx context.Context, _ *cli.Command) error {
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
		r.Use(middleware.CSP(csp))

		// Pages rendering
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
		r.Get("/*", web.RenderFunc(q, pool, publicURL))
		r.Handle("/static/*", web.StaticFunc())

		log.Info().Str("listenAddress", listenAddress).Msg("listening")
		return http.ListenAndServe(listenAddress, r)
	},
}

func main() {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")
	log.Logger = log.Logger.With().Caller().Logger()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("app crashed")
	}
}
