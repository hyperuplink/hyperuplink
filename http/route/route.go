package route

import (
	"strings"
)

// TODO: Find a way to build this automatically
var Routes map[string]Route = map[string]Route{
	"Account":                 {"account"},
	"AccountPassword":         {"account", "password"},
	"AccountProfile":          {"account", "profile"},
	"AccountSettings":         {"account", "settings"},
	"AccountTwofactor":        {"account", "twofactor"},
	"Admin":                   {"admin"},
	"AdminAuth":               {"admin", "auth"},
	"AdminBoardAttachments":   {"admin", "board", "attachments"},
	"AdminBoardCategories":    {"admin", "board", "categories"},
	"AdminBoardForums":        {"admin", "board", "forums"},
	"AdminBoardPosts":         {"admin", "board", "posts"},
	"AdminBoardProfiles":      {"admin", "board", "profiles"},
	"AdminBoardSignatures":    {"admin", "board", "signatures"},
	"AdminBoardTheme":         {"admin", "board", "theme"},
	"AdminCommsEmail":         {"admin", "comms", "email"},
	"AdminCommsXmpp":          {"admin", "comms", "xmpp"},
	"AdminGeneral":            {"admin", "general"},
	"AdminLog":                {"admin", "log"},
	"AdminUsers":              {"admin", "users"},
	"Categories":              {"_:categories"},
	"CategoriesForums":        {"_:categories", ":forums"},
	"CategoriesForumsTopics":  {"_:categories", ":forums", ":topics"},
	"DocsAbout":               {"docs", "about"},
	"DocsContact":             {"docs", "contact"},
	"DocsManual":              {"docs", "manual"},
	"DocsPrivacy":             {"docs", "privacy"},
	"DocsTerms":               {"docs", "terms"},
	"Session":                 {"session"},
	"SessionConfirm":          {"session", "confirm"},
	"SessionConfirmResend":    {"session", "confirm", "resend"},
	"SessionSettings":         {"session", "settings"},
	"SessionSignin":           {"session", "signin"},
	"SessionSignout":          {"session", "signout"},
	"SessionSignup":           {"session", "signup"},
	"SessionProvider":         {"session", ":provider"},
	"SessionProviderCallback": {"session", ":provider", "callback"},
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
