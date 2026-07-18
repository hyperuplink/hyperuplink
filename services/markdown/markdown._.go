package markdown

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Markdown struct {
	md   goldmark.Markdown
	docs goldmark.Markdown
}

func extensions() goldmark.Option {
	return goldmark.WithExtensions(
		extension.Strikethrough,
		extension.TaskList,
		extension.NewTable(),
		highlighting.NewHighlighting(
			// https://xyproto.github.io/splash/docs/
			highlighting.WithStyle("xcode"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
			),
		),
	)
}

func New() (*Markdown, error) {
	md := new(Markdown)

	md.md = goldmark.New(
		extensions(),
	)

	md.docs = goldmark.New(
		extensions(),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return md, nil
}

func (md *Markdown) Startup() (err error) {
	return nil
}

func (md *Markdown) Shutdown() (err error) {
	return nil
}

func convert(md goldmark.Markdown, src string) (dst string, err error) {
	var buf bytes.Buffer
	if err = md.Convert([]byte(src), &buf); err != nil {
		return dst, err
	}

	dst = buf.String()
	return dst, nil
}

func (md *Markdown) Convert(src string) (dst string, err error) {
	return convert(md.md, src)
}

func (md *Markdown) ConvertDocs(src string) (dst string, err error) {
	return convert(md.docs, src)
}
