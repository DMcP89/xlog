package mathjax

import (
	. "github.com/DMcP89/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	RegisterExtension(Mathjax{})
}

type Mathjax struct{}

func (Mathjax) Name() string { return "mathjax" }
func (Mathjax) Init() {
	RegisterStaticDir(js)
	registerBuildFiles()
	MarkdownConverter().Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&inlineMathParser{}, 999),
		),
		parser.WithBlockParsers(
			util.Prioritized(&mathJaxBlockParser{}, 999),
		),
	)
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&InlineMathRenderer{startDelim: `\(`, endDelim: `\)`}, 0),
		util.Prioritized(&MathBlockRenderer{startDelim: `\[`, endDelim: `\]`}, 0),
	))
}
