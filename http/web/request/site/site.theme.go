package site

import (
	"path"
	"strings"

	"github.com/mrusme/hyperuplink/models/setting"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

func (s *Site) themeSetting() (theme *setting.Theme, ok bool) {
	settingTheme, err := settingRepo.GetByID[setting.Theme](
		s.r.GetRuntime().Repositories.Setting,
		"theme",
	)
	if err != nil {
		return nil, false
	}

	return &settingTheme.JSONValue, true
}

func (s *Site) customBannerURL() (dlurl string, abs bool, ok bool) {
	theme, ok := s.themeSetting()
	if !ok || theme.CustomBanner == "" || theme.ThemeStorageProviderID == "" {
		return "", false, false
	}

	var err error
	if dlurl, abs, err = s.r.GetRuntime().Storage.GetFileDownloadURL(
		theme.ThemeStorageProviderID,
		path.Join(theme.ThemeStoragePath, theme.CustomBanner),
	); err != nil {
		return "", false, false
	}

	return dlurl, abs, true
}

func (s *Site) CustomBanner() (dlurl string) {
	dlurl, abs, ok := s.customBannerURL()
	if !ok {
		return ""
	}

	if !abs {
		dlurl = s.HrefTo(strings.TrimPrefix(dlurl, "/"))
	}

	return dlurl
}

func (s *Site) customFaviconURL() (dlurl string, abs bool, ok bool) {
	theme, ok := s.themeSetting()
	if !ok || theme.CustomFavicon == "" || theme.ThemeStorageProviderID == "" {
		return "", false, false
	}

	var err error
	if dlurl, abs, err = s.r.GetRuntime().Storage.GetFileDownloadURL(
		theme.ThemeStorageProviderID,
		path.Join(theme.ThemeStoragePath, theme.CustomFavicon),
	); err != nil {
		return "", false, false
	}

	return dlurl, abs, true
}

func (s *Site) CustomFavicon() (dlurl string) {
	dlurl, abs, ok := s.customFaviconURL()
	if !ok {
		return ""
	}

	if !abs {
		dlurl = s.HrefTo(strings.TrimPrefix(dlurl, "/"))
	}

	return dlurl
}

func (s *Site) BannerImageURL() string {
	if banner := s.CustomBanner(); banner != "" {
		return banner
	}

	return s.StaticFile("images/banner.jpg")
}

func (s *Site) CanonicalURL(baseURL string) string {
	return strings.TrimRight(baseURL, "/") + s.c.Path()
}

func (s *Site) OGImageURL(baseURL string) string {
	base := strings.TrimRight(baseURL, "/")

	if dlurl, abs, ok := s.customBannerURL(); ok {
		if abs {
			return dlurl
		}
		return base + "/" + strings.TrimPrefix(dlurl, "/")
	}

	return base + "/static/images/banner.jpg"
}
