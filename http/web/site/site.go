package site

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
)

type Site struct {
	rt      *runtime.Runtime
	relRoot string
	title   string
}

func New(rt *runtime.Runtime, c fiber.Ctx, title string) *Site {
	p := new(Site)

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

func (p *Site) GetRelRoot() string {
	return p.relRoot
}

func (p *Site) HrefTo(path string) string {
	return fmt.Sprintf("%s%s", p.GetRelRoot(), path)
}

func (p *Site) StaticFile(filename string) string {
	hash := p.rt.Build.Hash

	if p.rt.IsDevelopmentMode() {
		hash = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	return fmt.Sprintf("%sstatic/%s?v=%s",
		p.GetRelRoot(),
		filename,
		hash,
	)
}

func (p *Site) CSS(name string) string {
	return p.StaticFile("css/" + name)
}

func (p *Site) Title() string {
	return p.title
}
