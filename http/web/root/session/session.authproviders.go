package session

import (
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/reddit"
	"github.com/markbates/goth/providers/slack"

	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
)

func authProviderCallbackURL(baseURL, name string) string {
	return strings.TrimRight(baseURL, "/") + "/session/" + name + "/callback"
}

func buildAuthProvider(ap config.AuthProvider, baseURL string) goth.Provider {
	name := strings.ToLower(ap.Type)
	callbackURL := authProviderCallbackURL(baseURL, name)

	switch name {
	case "github":
		return github.New(ap.Key, ap.Secret, callbackURL, ap.Scopes...)
	case "google":
		return google.New(ap.Key, ap.Secret, callbackURL, ap.Scopes...)
	case "facebook":
		return facebook.New(ap.Key, ap.Secret, callbackURL, ap.Scopes...)
	case "slack":
		return slack.New(ap.Key, ap.Secret, callbackURL, ap.Scopes...)
	case "reddit":
		p := reddit.New(ap.Key, ap.Secret, callbackURL, "temporary", "", "", ap.Scopes...)
		return &p
	default:
		return nil
	}
}

func RegisterAuthProviders(rt *runtime.Runtime, baseURL string) error {
	providers, err := rt.Config.AuthProviders()
	if err != nil {
		return err
	}

	var built []goth.Provider
	for _, ap := range providers {
		if ap.Key == "" || ap.Secret == "" {
			continue
		}

		provider := buildAuthProvider(ap, baseURL)
		if provider == nil {
			rt.Warn("auth_provider", "unsupported provider type", "type", ap.Type)
			continue
		}

		built = append(built, provider)
	}

	goth.UseProviders(built...)

	return nil
}
