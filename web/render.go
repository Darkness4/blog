//go:generate go run -tags build build.go
package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Darkness4/blog/db"
	"github.com/Darkness4/blog/utils/color"
	"github.com/Darkness4/blog/utils/math"
	"github.com/Darkness4/blog/web/gen/index"
	"github.com/Masterminds/sprig/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed gen components base.html base.htmx error.tmpl
	html embed.FS
	//go:embed static
	static embed.FS
)

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

var funcsMap = func() template.FuncMap {
	f := sprig.TxtFuncMap()
	f["computeColorByWord"] = color.ComputeByWord
	return f
}

func RenderFunc(
	q *db.Queries,
	pool *pgxpool.Pool,
	publicURL string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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
		} else if !errors.Is(err, fs.ErrNotExist) {
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
		if err != nil {
			if strings.Contains(err.Error(), "no files") {
				// Render 404
				if err := renderError(
					w,
					r,
					html,
					"Oops! The page you were looking for couldn't be found.",
					http.StatusNotFound,
				); err != nil {
					log.Err(err).Msg("failed to render error")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.Err(err).Msg("template error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageS := r.URL.Query().Get("page")
		page, _ := strconv.Atoi(pageS)

		rctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pv, err := q.FindPageViewsOrZero(rctx, strings.ToLower(cleanPath))
		if err != nil {
			log.Err(err).Msg("failed to fetch page views")
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
			log.Err(err).Msg("failed to execute template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cleanPath != "/" {
			go func() {
				rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := q.CreateOrIncrementPageViewsOnUniqueIP(rctx, pool, strings.ToLower(cleanPath), ReadUserIP(r)); err != nil {
					log.Err(err).Msg("failed to increment page views")
				}
			}()
		}
	}
}

// renderError renders an error page on error.
func renderError(
	w http.ResponseWriter, r *http.Request,
	html embed.FS,
	errorMsg string,
	code int,
) error {
	var base string
	if r.Header.Get("Hx-Boosted") != "true" {
		// Initial Rendering
		base = "base.html"
	} else {
		// SSR
		base = "base.htmx"
	}
	ctx := r.Context()
	ctx = context.WithValue(ctx, "error_code", code)
	ctx = context.WithValue(ctx, "error_short", http.StatusText(code))
	ctx = context.WithValue(
		ctx,
		"error_long",
		errorMsg,
	)

	t, err := template.New("base").
		Funcs(funcsMap()).
		ParseFS(html, base, "error.tmpl", "components/*.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse error.tmpl: %v", err))
	}
	return t.ExecuteTemplate(w, "base", struct {
		Context context.Context
	}{
		Context: ctx,
	})
}

func StaticFunc() http.Handler {
	return http.FileServer(http.FS(static))
}
