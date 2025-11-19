package page

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
)

type Page struct {
	rt      *runtime.Runtime
	relRoot string
	title   string
}

func New(rt *runtime.Runtime, c fiber.Ctx, title string) *Page {
	p := new(Page)

	p.rt = rt

	cR := c.Route()
	parts := strings.Count(cR.Path, "/")

	relRoot := ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}
	p.relRoot = relRoot

	p.title = title

	return p
}

func (p *Page) StaticFile(filename string) string {
	hash := p.rt.Build.Hash

	if p.rt.IsDevelopmentMode() {
		hash = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	return fmt.Sprintf("%sstatic/%s?v=%s",
		p.relRoot,
		filename,
		hash,
	)
}

func (p *Page) CSS(name string) string {
	return p.StaticFile("css/" + name)
}

func (p *Page) Title() string {
	return p.title
}
