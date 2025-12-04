package route

import (
	"strings"
)

// TODO: Find a way to build this automatically
var Routes map[string]Route = map[string]Route{
	"Account":                 {"account"},
	"AccountSettings":         {"account", "settings"},
	"AccountProfile":          {"account", "profile"},
	"Session":                 {"session"},
	"SessionConfirm":          {"session", "confirm"},
	"SessionConfirmResend":    {"session", "confirm", "resend"},
	"SessionSignup":           {"session", "signun"},
	"SessionSignin":           {"session", "signin"},
	"SessionSignout":          {"session", "signout"},
	"SessionProvider":         {"session", ":provider"},
	"SessionProviderCallback": {"session", ":provider", "callback"},
	"Admin":                   {"admin"},
	"Categories":              {":categories"},
	"CategoriesForums":        {":categories", ":forums"},
	"CategoriesForumsTopics":  {":categories", ":forums", ":topics"},
	"System":                  {"system"},
}

type Route []string

func (r Route) AsURL() string {
	return strings.Join([]string(r), "/")
}

func (r Route) Pathname() string {
	rl := len(r)
	if rl == 0 {
		return ""
	}

	return r[rl-1]
}

func For(id string) (r Route) {
	var ok bool

	if r, ok = Routes[id]; !ok {
		return Route{}
	}

	return r
}
