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

// Replacer render image with replaced source link.
type Replacer struct {
	html.Config
	ReplaceFunc
}

// New return initialized image render with source url replacing support.
func New(r ReplaceFunc, options ...html.Option) *Replacer {
	var config = html.NewConfig()
	for _, opt := range options {
		opt.SetHTMLOption(&config)
	}
	return &Replacer{
		Config:      config,
		ReplaceFunc: r,
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs interface.
func (r *Replacer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *Replacer) renderImage(
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

	if _, err := w.WriteString("<img src=\""); err != nil {
		return ast.WalkSkipChildren, err
	}
	if _, err := w.Write(n.Destination); err != nil {
		return ast.WalkSkipChildren, err
	}
	if _, err := w.WriteString(`" alt="`); err != nil {
		return ast.WalkSkipChildren, err
	}
	if _, err := w.Write(n.Text(source)); err != nil {
		return ast.WalkSkipChildren, err
	}
	if err := w.WriteByte('"'); err != nil {
		return ast.WalkSkipChildren, err
	}
	if n.Title != nil {
		if _, err := w.WriteString(` title="`); err != nil {
			return ast.WalkSkipChildren, err
		}
		r.Writer.Write(w, n.Title)
		if err := w.WriteByte('"'); err != nil {
			return ast.WalkSkipChildren, err
		}
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		if _, err := w.WriteString(" />"); err != nil {
			return ast.WalkSkipChildren, err
		}
	} else {
		if _, err := w.WriteString(">"); err != nil {
			return ast.WalkSkipChildren, err
		}
	}
	return ast.WalkSkipChildren, nil
}

func (r *Replacer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	if n.IsCode() {
		_, _ = w.Write(n.Value)
	} else {
		if n.IsRaw() {
			r.Writer.RawWrite(w, n.Value)
		} else {
			r.Writer.Write(w, n.Value)
		}
	}
	return ast.WalkContinue, nil
}

// Extend implement goldmark.Extender interface.
func (r *Replacer) Extend(m goldmark.Markdown) {
	if r.ReplaceFunc == nil {
		return
	}
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(r, 0),
		),
	)
}
