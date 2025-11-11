package search

import (
	"html/template"
	"net/http"

	_ "embed"

	"github.com/Darkness4/blog/meilisearch"
	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

//go:embed search.tpl
var searchTemplate string

func recordsGroupByLvl0(
	records []meilisearch.RecordWithFormat,
) (m map[string][]meilisearch.RecordWithFormat, keys []string) {
	m = make(map[string][]meilisearch.RecordWithFormat)

	for _, r := range records {
		if _, ok := m[r.HierarchyLvl0]; !ok {
			keys = append(keys, r.HierarchyLvl0)
			m[r.HierarchyLvl0] = make([]meilisearch.RecordWithFormat, 0, 1)
		}
		m[r.HierarchyLvl0] = append(m[r.HierarchyLvl0], r)
	}
	return
}

func funcsMap() template.FuncMap {
	m := sprig.HtmlFuncMap()
	m["noescape"] = func(s string) template.HTML { return template.HTML(s) }
	return m
}

func Handler(meili *meilisearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query().Get("q")

		if q == "" {
			return
		}

		res, err := meili.Search(ctx, q)
		if err != nil {
			log.Err(err).Msg("search failure")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var records []meilisearch.RecordWithFormat
		if err := res.Hits.DecodeInto(&records); err != nil {
			log.Err(err).Msg("search failure")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		groupedRecords, lvl0s := recordsGroupByLvl0(records)

		if err := template.Must(template.New("base").
			Funcs(funcsMap()).
			Parse(searchTemplate)).
			Execute(w, map[string]any{
				"Lvl0s":          lvl0s,
				"GroupedRecords": groupedRecords,
			}); err != nil {
			log.Err(err).Msg("template error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}
}
