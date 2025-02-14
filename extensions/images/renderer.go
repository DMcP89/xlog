package images

import (
	. "github.com/DMcP89/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type imagesColumnsRenderer struct{}

func (s *imagesColumnsRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindColumns, s.render)
}

func (s *imagesColumnsRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(`<div class="columns">`)

		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			w.WriteString(`<div class="column">`)
			MarkdownConverter().Renderer().Render(w, source, c)
			w.WriteString(`</div>`)
		}

	} else {
		w.WriteString(`</div>`)
	}

	return ast.WalkSkipChildren, nil
}
