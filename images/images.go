// Package images is a extension for the goldmark (http://github.com/yuin/goldmark).
//
// This extension adds replacer render to change image urls.
package images

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// ReplaceFunc is a function for replacing image source link.
type ReplaceFunc = func(link string) string

// NewReplacer adding src url replacing function to image html render.
func NewReplacer(r ReplaceFunc) goldmark.Option {
	return goldmark.WithRendererOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(New(r), 0),
		),
	)
}

// replacer render image with replaced source link.
type replacer struct {
	html.Config
	ReplaceFunc
}

// New return initialized image render with source url replacing support.
func New(r ReplaceFunc, options ...html.Option) goldmark.Extender {
	var config = html.NewConfig()
	for _, opt := range options {
		opt.SetHTMLOption(&config)
	}
	return &replacer{
		Config:      config,
		ReplaceFunc: r,
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs interface.
func (r *replacer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *replacer) renderImage(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	// add image link replacing hack
	if r.ReplaceFunc != nil {
		var src = r.ReplaceFunc(util.BytesToReadOnlyString(n.Destination))
		// if src == "" {
		// 	return ast.WalkContinue, nil
		// } else if src == "-" {
		// 	return ast.WalkSkipChildren, nil
		// }
		n.Destination = util.StringToReadOnlyBytes(src)
	}

	w.WriteString("<img src=\"")
	w.Write(n.Destination)
	w.WriteString(`" alt="`)
	w.Write(n.Text(source))
	w.WriteByte('"')
	if n.Title != nil {
		w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		w.WriteString(" />")
	} else {
		w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

// Extend implement goldmark.Extender interface.
func (r *replacer) Extend(m goldmark.Markdown) {
	if r.ReplaceFunc == nil {
		return
	}
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(r, 0),
		),
	)
}
