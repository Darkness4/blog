package d2

import (
	"bytes"
	"context"

	"cdr.dev/slog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

type CompileOptions d2lib.CompileOptions
type RenderOptions d2svg.RenderOpts

type HTMLRenderer struct {
	CompileOptions
	RenderOptions
}

func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindBlock, r.Render)
}

func (r *HTMLRenderer) Render(
	w util.BufWriter,
	src []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*Block)
	if !entering {
		_, _ = w.WriteString("</div>")
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<div class="d2">`)

	b := bytes.Buffer{}
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		b.Write(line.Value(src))
	}

	if b.Len() == 0 {
		return ast.WalkContinue, nil
	}

	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return ast.WalkStop, err
	}
	compileOpts := d2lib.CompileOptions(r.CompileOptions)
	if compileOpts.Ruler == nil {
		compileOpts.Ruler = ruler
	}
	if compileOpts.LayoutResolver == nil {
		compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
			return d2dagrelayout.DefaultLayout, nil
		}
	}
	renderOpts := d2svg.RenderOpts(r.RenderOptions)
	diagram, _, err := d2lib.Compile(
		log.With(context.Background(), slog.Make()),
		b.String(),
		&compileOpts,
		&renderOpts,
	)
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}
	out, err := d2svg.Render(diagram, &renderOpts)
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}

	_, err = w.Write(out)
	return ast.WalkContinue, err
}
