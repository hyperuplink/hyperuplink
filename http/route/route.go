package route

import (
	"strings"
)

type Route struct {
	hierarchy    []string
	noBreadcrumb bool
}

// TODO: Find a way to build this automatically
var Routes map[string]Route = map[string]Route{
	"Root":                      {hierarchy: []string{"root"}},
	"Account":                   {hierarchy: []string{"account"}},
	"AccountPassword":           {hierarchy: []string{"account", "password"}},
	"AccountProfile":            {hierarchy: []string{"account", "profile"}},
	"AccountSettings":           {hierarchy: []string{"account", "settings"}},
	"AccountTwofactor":          {hierarchy: []string{"account", "twofactor"}},
	"Admin":                     {hierarchy: []string{"admin"}},
	"AdminAuth":                 {hierarchy: []string{"admin", "auth"}},
	"AdminBoard":                {hierarchy: []string{"admin", "board"}},
	"AdminBoardAttachments":     {hierarchy: []string{"admin", "board", "attachments"}},
	"AdminBoardCategories":      {hierarchy: []string{"admin", "board", "categories"}},
	"AdminBoardForums":          {hierarchy: []string{"admin", "board", "forums"}},
	"AdminBoardTopics":          {hierarchy: []string{"admin", "board", "topics"}},
	"AdminBoardProfiles":        {hierarchy: []string{"admin", "board", "profiles"}},
	"AdminBoardSignatures":      {hierarchy: []string{"admin", "board", "signatures"}},
	"AdminBoardTheme":           {hierarchy: []string{"admin", "board", "theme"}},
	"AdminCommsEmail":           {hierarchy: []string{"admin", "comms", "email"}},
	"AdminCommsXmpp":            {hierarchy: []string{"admin", "comms", "xmpp"}},
	"AdminGeneral":              {hierarchy: []string{"admin", "general"}},
	"AdminLog":                  {hierarchy: []string{"admin", "log"}},
	"AdminUsers":                {hierarchy: []string{"admin", "users"}},
	"Categories":                {hierarchy: []string{"_:categories"}},
	"CategoriesForums":          {hierarchy: []string{"_:categories", ":forums"}},
	"CategoriesForumsTopics":    {hierarchy: []string{"_:categories", ":forums", ":topics"}},
	"CategoriesForumsTopicsNew": {hierarchy: []string{"_:categories", ":forums", ":topics", "new"}},
	"Docs":                      {hierarchy: []string{"docs"}, noBreadcrumb: true},
	"DocsAbout":                 {hierarchy: []string{"docs", "about"}},
	"DocsContact":               {hierarchy: []string{"docs", "contact"}},
	"DocsManual":                {hierarchy: []string{"docs", "manual"}},
	"DocsPrivacy":               {hierarchy: []string{"docs", "privacy"}},
	"DocsTerms":                 {hierarchy: []string{"docs", "terms"}},
	"Search":                    {hierarchy: []string{"search"}},
	"Session":                   {hierarchy: []string{"session"}, noBreadcrumb: true},
	"SessionConfirm":            {hierarchy: []string{"session", "confirm"}},
	"SessionConfirmResend":      {hierarchy: []string{"session", "confirm", "resend"}},
	"SessionSettings":           {hierarchy: []string{"session", "settings"}},
	"SessionSignin":             {hierarchy: []string{"session", "signin"}},
	"SessionSignout":            {hierarchy: []string{"session", "signout"}},
	"SessionSignup":             {hierarchy: []string{"session", "signup"}},
	"SessionProvider":           {hierarchy: []string{"session", ":provider"}},
	"SessionProviderCallback":   {hierarchy: []string{"session", ":provider", "callback"}},
	"System":                    {hierarchy: []string{"system"}},
}

func (r Route) Len() int {
	return len(r.hierarchy)
}

func (r Route) AsURL() string {
	return strings.Join([]string(r.hierarchy), "/")
}

func (r Route) AsTitle() string {
	joined := strings.Join([]string(r.hierarchy), "_")
	if strings.Index(joined, ":") > -1 {
		return ""
	}
	return joined + "_title"
}

func (r Route) Pathname() string {
	rl := r.Len()
	if rl == 0 {
		return ""
	}

	return r.hierarchy[rl-1]
}

func (r Route) Parent() (id string) {
	rl := r.Len()
	if rl < 2 {
		return ""
	}

	parentRoute := Route{hierarchy: r.hierarchy[:rl-1]}
	for key, rt := range Routes {
		if rt.Equals(parentRoute) {
			return key
		}
	}

	return ""
}

func (r Route) Equals(cmp Route) (eq bool) {
	if cmp.Len() != r.Len() {
		return false
	}
	for idx := range r.hierarchy {
		if cmp.hierarchy[idx] != r.hierarchy[idx] {
			return false
		}
	}

	return true
}

func (r Route) ParentRoute() (rt Route) {
	id := r.Parent()
	if id == "" {
		return Route{}
	}

	return For(id)
}

func (r Route) HasBreadcrumb() bool {
	return r.noBreadcrumb == false
}

func For(id string) (r Route) {
	var ok bool

	if r, ok = Routes[id]; !ok {
		return Route{}
	}

	return r
}
