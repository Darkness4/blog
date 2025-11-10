// Package template provides template helpers.
package template

import (
	"strings"
	"text/template"
)

// Quick renders a template and panics on error.
func Quick(t string, data any) string {
	buf := strings.Builder{}
	tmpl := template.Must(template.New("template").Parse(t))
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

// Render renders a template.
func Render(t string, data any) (string, error) {
	buf := strings.Builder{}
	tmpl, err := template.New("template").Parse(t)
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
