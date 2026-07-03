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
	"github.com/mrusme/hyperuplink/models/setting"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

type Site struct {
	r        route.IRouteController
	c        fiber.Ctx
	csrf     string
	pathName string
	absPath  string
	relRoot  string
	pager    *Pager

	rt       route.Route
	title    string
	timezone *time.Location
}

func New(r route.IRouteController, c fiber.Ctx) *Site {
	s := new(Site)

	s.r = r
	s.c = c
	s.csrf = csrf.TokenFromContext(s.c)
	s.pager = NewPager(1, 1, 1)

	s.pathName, s.absPath, s.relRoot = helpers.GetPaths(s.c)

	s.SetTimezone("UTC")

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

func (s *Site) GetFormAction(args ...string) string {
	pathname := s.GetPathname()

	if len(args) == 0 {
		return pathname
	}

	target := args[0]

	if target == "" || target == pathname {
		return pathname
	}

	return pathname + "/" + target
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
	var staticPicture string = s.StaticFile("images/avatar.jpg")

	if id == "" {
		return staticPicture
	}

	var settingProfiles *setting.Setting[setting.Profiles]
	settingProfiles, err = settingRepo.GetByID[setting.Profiles](
		s.r.GetRuntime().Repositories.Setting,
		"profiles",
	)
	if err != nil {
		return staticPicture
	}
	profiles := settingProfiles.JSONValue

	if !profiles.EnablePicture || profiles.PictureStorageProviderID == "" {
		return staticPicture
	}

	if dlurl, err = s.r.GetRuntime().Storage.GetFileDownloadURL(
		profiles.PictureStorageProviderID,
		profiles.PictureStoragePath+"/"+id+"."+profiles.PictureFormat,
	); err != nil {
		return staticPicture
	}

	return dlurl
}

func (s *Site) SetRoute(rt route.Route) {
	s.rt = rt
}

func (s *Site) GetRoute() route.Route {
	return s.rt
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

func (s *Site) SetTimezone(tz string) (err error) {
	if s.timezone, err = time.LoadLocation(tz); err != nil {
		return err
	}

	return nil
}

func (s *Site) Date(ts pgtype.Timestamp) (date string) {
	if ts.Valid == false {
		return "-"
	}

	return ts.Time.In(s.timezone).Format("2006-01-02")
}

func (s *Site) DateTime(ts pgtype.Timestamp) (timedate string) {
	if ts.Valid == false {
		return "-"
	}

	return ts.Time.In(s.timezone).Format("2006-01-02, 15:04 MST")
}
