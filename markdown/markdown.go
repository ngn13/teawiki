package markdown

import (
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	bf "github.com/russross/blackfriday/v2"
)

// modified code from: https://github.com/depado/bfchroma/blob/main/renderer.go

type Markdown struct {
	Formatter  *html.Formatter
	Style      *chroma.Style
	Renderer   bf.Renderer
	Extensions bf.Extensions
}

func New(style string) *Markdown {
	md := &Markdown{
		Formatter: html.New(),
		Style:     styles.Get(style),
		Renderer: bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		}),
		Extensions: bf.NoIntraEmphasis | bf.Tables |
			bf.FencedCode | bf.DefinitionLists | bf.Footnotes |
			bf.Autolink | bf.HeadingIDs | bf.AutoHeadingIDs |
			bf.BackslashLineBreak,
	}

	return md
}

func (md *Markdown) RenderWithChroma(
	w io.Writer,
	text []byte,
	data bf.CodeBlockData,
) error {
	var (
		lexer    chroma.Lexer
		iterator chroma.Iterator
		err      error
	)

	if len(data.Info) > 0 {
		lexer = lexers.Get(string(data.Info))
	} else {
		lexer = lexers.Analyse(string(text))
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	if iterator, err = lexer.Tokenise(nil, string(text)); err != nil {
		return err
	}

	return md.Formatter.Format(w, md.Style, iterator)
}

func (md *Markdown) RenderNode(
	w io.Writer,
	node *bf.Node,
	entering bool,
) bf.WalkStatus {
	switch node.Type {
	case bf.Document:
		if entering {
			w.Write([]byte("<style>"))
			md.Formatter.WriteCSS(w, md.Style)
			w.Write([]byte("</style>"))
		}

		return md.Renderer.RenderNode(w, node, entering)

	case bf.CodeBlock:
		if err := md.RenderWithChroma(
			w, node.Literal, node.CodeBlockData,
		); err != nil {
			return md.Renderer.RenderNode(w, node, entering)
		}

		return bf.SkipChildren
	}

	return md.Renderer.RenderNode(w, node, entering)
}

func (md *Markdown) RenderHeader(w io.Writer, ast *bf.Node) {
	md.Renderer.RenderHeader(w, ast)
}

func (md *Markdown) RenderFooter(w io.Writer, ast *bf.Node) {
	md.Renderer.RenderFooter(w, ast)
}

func (md *Markdown) Render(r io.Reader) []byte {
	if content, err := io.ReadAll(r); err != nil {
		return nil
	} else {
		return bf.Run(
			content, bf.WithRenderer(md), bf.WithExtensions(md.Extensions))
	}
}
