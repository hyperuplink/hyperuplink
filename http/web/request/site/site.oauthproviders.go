package site

import (
	"sort"
	"strings"

	"github.com/markbates/goth"
)

type OauthProvider struct {
	Label string
	Title string
	Href  string
	Class string
}

var oauthProviderLabels = map[string]string{
	"github":   "GitHub",
	"google":   "Google",
	"facebook": "Facebook",
	"reddit":   "Reddit",
	"slack":    "Slack",
	"apple":    "Apple",
	"twitter":  "X",
	"bsky":     "Bluesky",
}

func oauthProviderLabel(name string) string {
	if name == "" {
		return ""
	}

	if label, ok := oauthProviderLabels[name]; ok {
		return label
	}

	return strings.ToUpper(name[:1]) + name[1:]
}

func (s *Site) OauthProviderLabel(name string) string {
	return oauthProviderLabel(strings.ToLower(name))
}

func (s *Site) OauthProviders() []OauthProvider {
	registered := goth.GetProviders()

	names := make([]string, 0, len(registered))
	for name := range registered {
		names = append(names, name)
	}
	sort.Strings(names)

	providers := make([]OauthProvider, 0, len(names))
	for _, name := range names {
		label := oauthProviderLabel(name)
		providers = append(providers, OauthProvider{
			Label: label,
			Title: label,
			Href:  s.HrefRoute("session", name),
			Class: name,
		})
	}

	return providers
}
