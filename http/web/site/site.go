package site

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/kaptinlin/go-i18n"
	"github.com/mrusme/hyperuplink/http/route"
)

type Site struct {
	r       route.IRoute
	c       fiber.Ctx
	csrf    string
	absPath string
	relRoot string
	title   string

	I18n *i18n.Localizer
}

func New(r route.IRoute, c fiber.Ctx) *Site {
	s := new(Site)

	s.r = r
	s.c = c
	s.csrf = csrf.TokenFromContext(s.c)

	cR := s.c.Route()
	s.absPath = cR.Path
	parts := strings.Count(s.absPath, "/")

	relRoot := ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}
	s.relRoot = relRoot

	acceptLang := c.Get("Accept-Language", "en")
	s.I18n = r.GetRuntime().Intnat.NewLocalizer(acceptLang)

	return s
}

func (s *Site) T(msg string) string {
	return s.I18n.Get(msg)
}

func (s *Site) GetRelRoot() string {
	return s.relRoot
}

func (s *Site) GetAbsPath() string {
	return s.absPath
}

func (s *Site) GetCSRFToken() string {
	return s.csrf
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
