package d2

import (
	"bytes"
	"context"
	"log/slog"
	"os"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
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

type immutableAttributes struct {
	n ast.Node
}

func (a *immutableAttributes) Get(name []byte) (any, bool) {
	return a.n.Attribute(name)
}

func (a *immutableAttributes) GetString(name string) (any, bool) {
	return a.n.AttributeString(name)
}

func (a *immutableAttributes) All() []ast.Attribute {
	if a.n.Attributes() == nil {
		return []ast.Attribute{}
	}
	return a.n.Attributes()
}

func getAttributes(infostr []byte) *immutableAttributes {
	if infostr != nil {
		attrStartIdx := -1

		for idx, char := range infostr {
			if char == '{' {
				attrStartIdx = idx
				break
			}
		}
		if attrStartIdx > 0 {
			n := ast.NewTextBlock() // dummy node for storing attributes
			attrStr := infostr[attrStartIdx:]
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				for _, attr := range attrs {
					n.SetAttribute(attr.Name, attr.Value)
				}
				return &immutableAttributes{n}
			}
		}
	}
	return nil
}

func (r *HTMLRenderer) Render(
	w util.BufWriter,
	src []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*Block)
	if !entering {
		_, _ = w.WriteString("</figure></center></div>")
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<div class="d2"><center><figure>`)

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
	// Parsing Info ```d2 ({key=value}) <-- THIS
	var title string
	if attrs := getAttributes(n.Info()); attrs != nil {
		if layoutAttr, ok := attrs.GetString("layout"); ok {
			if layoutAttrB, ok := layoutAttr.([]uint8); ok {
				layoutAttrStr := string(layoutAttrB)
				switch layoutAttrStr {
				case "dagre":
					compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
						return d2dagrelayout.DefaultLayout, nil
					}
				case "elk":
					compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
						return d2elklayout.DefaultLayout, nil
					}
				}
			}
		}
		if titleAttr, ok := attrs.GetString("title"); ok {
			if titleAttrB, ok := titleAttr.([]uint8); ok {
				title = string(titleAttrB)
			}
		}
	}

	if compileOpts.LayoutResolver == nil {
		compileOpts.LayoutResolver = func(_ string) (d2graph.LayoutGraph, error) {
			return d2dagrelayout.DefaultLayout, nil
		}
	}
	renderOpts := d2svg.RenderOpts(r.RenderOptions)
	diagram, _, err := d2lib.Compile(
		log.With(context.Background(), slog.New(slog.NewTextHandler(os.Stderr, nil))),
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
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}
	_, err = w.WriteString("<figcaption><i>")
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}
	_, err = w.WriteString(title)
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}
	_, err = w.WriteString("</i></figcaption>")
	if err != nil {
		_, _ = w.Write(b.Bytes())
		return ast.WalkContinue, err
	}
	return ast.WalkContinue, err
}
