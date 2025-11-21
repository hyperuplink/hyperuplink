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
	s := new(Site)

	s.rt = rt

	cR := c.Route()
	parts := strings.Count(cR.Path, "/")

	relRoot := ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}
	s.relRoot = relRoot

	s.title = title

	return s
}

func (s *Site) GetRelRoot() string {
	return s.relRoot
}

func (s *Site) HrefTo(path string) string {
	return fmt.Sprintf("%s%s", s.GetRelRoot(), path)
}

func (s *Site) StaticFile(filename string) string {
	hash := s.rt.Build.Hash

	if s.rt.IsDevelopmentMode() {
		hash = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	return fmt.Sprintf("%sstatic/%s?v=%s",
		s.GetRelRoot(),
		filename,
		hash,
	)
}

func (s *Site) CSS(name string) string {
	return s.StaticFile("css/" + name)
}

func (s *Site) Title() string {
	return s.title
}
