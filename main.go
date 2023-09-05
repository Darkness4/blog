//go:generate go run build.go

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed gen components base.html base.htmx
	html embed.FS
	//go:embed static
	static        embed.FS
	version       = "dev"
	listenAddress string
)

var app = &cli.App{
	Name:    "blog",
	Version: version,
	Usage:   "A blog in HTMX.",
	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "listen.address",
			Value:       ":3000",
			Destination: &listenAddress,
		},
	},
	Action: func(cCtx *cli.Context) error {
		log.Level(zerolog.DebugLevel)

		// Router
		r := chi.NewRouter()
		r.Use(hlog.NewHandler(log.Logger))

		// Pages rendering
		var renderFn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Clean(r.URL.Path)
			path = filepath.Clean(fmt.Sprintf("gen/pages/%s/page.tmpl", path))

			// Check if SSR
			var base string
			if r.Header.Get("Hx-Request") != "true" {
				// Initial Rendering
				base = "base.html"
			} else {
				// SSR
				base = "base.htmx"
			}
			t, err := template.ParseFS(html, base, path, "components/*")
			if err != nil {
				// The page doesn't exist
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			if err := t.ExecuteTemplate(w, "base", struct{}{}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
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
