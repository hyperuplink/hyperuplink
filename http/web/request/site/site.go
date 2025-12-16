package site

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/helpers"
	"github.com/mrusme/hyperuplink/models/user"
)

type Site struct {
	r        route.IRouteController
	c        fiber.Ctx
	csrf     string
	pathName string
	absPath  string
	relRoot  string
	pager    *Pager

	title       string
	currentUser *user.User
}

func New(r route.IRouteController, c fiber.Ctx) *Site {
	s := new(Site)

	s.r = r
	s.c = c
	s.csrf = csrf.TokenFromContext(s.c)
	s.pager = NewPager(1, 1, 1)

	s.pathName, s.absPath, s.relRoot = helpers.GetPaths(s.c)

	return s
}

func (s *Site) GetRelRoot() string {
	return s.relRoot
}

func (s *Site) GetAbsPath() string {
	return s.absPath
}

func (s *Site) GetPathname() string {
	cidx := strings.Index(s.pathName, ":")
	if cidx == -1 {
		return s.pathName
	}

	segment := s.pathName[(cidx + 1):]

	r := route.Route{}
	r.SetHierarchy([]string{"root", s.pathName})
	rf := r.Fill(
		map[string]string{
			segment: s.c.Params(segment),
		},
	)

	return rf.AsURL()
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

func (s *Site) ProfilePicture(id string) (dlurl string) {
	var err error

	// TODO: Get provider ID from System
	// TODO: Get path from System
	// TODO: Get .webp from format configured in System
	if dlurl, err = s.r.GetRuntime().Storage.GetFileDownloadURL(
		"profile-pictures",
		"profile-pictures/"+id+".webp",
	); err != nil {
		return s.StaticFile("images/avatar.jpg")
	}

	return dlurl
}

func (s *Site) Title() string {
	return s.r.GetEnv().Title
}

func (s *Site) TitleFull(systemTitle string) string {
	return fmt.Sprintf("%s // %s",
		s.r.GetEnv().Title,
		systemTitle,
	)
}

func (s *Site) SetTitle(title string) {
	s.r.GetEnv().Title = title
}

func (s *Site) Date(ts pgtype.Timestamp) (date string) {
	if ts.Valid == false {
		return "-"
	}

	// TODO: Make it user configurable
	return ts.Time.Format("2006-01-02")
}

func (s *Site) DateTime(ts pgtype.Timestamp) (timedate string) {
	if ts.Valid == false {
		return "-"
	}

	// TODO: Make it user configurable
	// TODO: Handle user timezone
	return ts.Time.Format("2006-01-02, 15:04 MST")
}
