package site

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
)

type Site struct {
	r       route.IRoute
	relRoot string
	title   string
}

func New(r route.IRoute, c fiber.Ctx) *Site {
	s := new(Site)

	s.r = r

	cR := c.Route()
	parts := strings.Count(cR.Path, "/")

	relRoot := ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}
	s.relRoot = relRoot

	return s
}

func (s *Site) GetRelRoot() string {
	return s.relRoot
}

func (s *Site) HrefTo(path string) string {
	return fmt.Sprintf("%s%s", s.GetRelRoot(), path)
}

func (s *Site) HrefRoute(routes ...string) string {
	return s.HrefTo(strings.Join(routes, "/"))
}

func (s *Site) StaticFile(filename string) string {
	hash := s.r.GetRuntime().Build.Hash

	if s.r.GetRuntime().IsDevelopmentMode() {
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
	return s.r.GetEnv().Title
}
