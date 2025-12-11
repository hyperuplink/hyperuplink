package markdown

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Markdown struct {
	md goldmark.Markdown
}

func New() (*Markdown, error) {
	md := new(Markdown)

	md.md = goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				// https://xyproto.github.io/splash/docs/
				highlighting.WithStyle("xcode"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
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

func (md *Markdown) Convert(src string) (dst string, err error) {
	var buf bytes.Buffer
	if err = md.md.Convert([]byte(src), &buf); err != nil {
		return dst, err
	}

	dst = buf.String()
	return dst, nil
}
