package presentme

import (
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func intoHTML(w io.Writer, bytes []byte) error {
	return md.Convert(bytes, w)
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Typographer,
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
		),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)
