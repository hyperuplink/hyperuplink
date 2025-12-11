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
	"Root":                    {hierarchy: []string{"root"}},
	"Account":                 {hierarchy: []string{"root", "account"}},
	"AccountPassword":         {hierarchy: []string{"root", "account", "password"}},
	"AccountProfile":          {hierarchy: []string{"root", "account", "profile"}},
	"AccountSettings":         {hierarchy: []string{"root", "account", "settings"}},
	"AccountTwofactor":        {hierarchy: []string{"root", "account", "twofactor"}},
	"Admin":                   {hierarchy: []string{"root", "admin"}},
	"AdminAuth":               {hierarchy: []string{"root", "admin", "auth"}},
	"AdminBoard":              {hierarchy: []string{"root", "admin", "board"}},
	"AdminBoardAttachments":   {hierarchy: []string{"root", "admin", "board", "attachments"}},
	"AdminBoardCategories":    {hierarchy: []string{"root", "admin", "board", "categories"}},
	"AdminBoardForums":        {hierarchy: []string{"root", "admin", "board", "forums"}},
	"AdminBoardTopics":        {hierarchy: []string{"root", "admin", "board", "topics"}},
	"AdminBoardProfiles":      {hierarchy: []string{"root", "admin", "board", "profiles"}},
	"AdminBoardSignatures":    {hierarchy: []string{"root", "admin", "board", "signatures"}},
	"AdminBoardTheme":         {hierarchy: []string{"root", "admin", "board", "theme"}},
	"AdminCommsEmail":         {hierarchy: []string{"root", "admin", "comms", "email"}},
	"AdminCommsXmpp":          {hierarchy: []string{"root", "admin", "comms", "xmpp"}},
	"AdminGeneral":            {hierarchy: []string{"root", "admin", "general"}},
	"AdminLog":                {hierarchy: []string{"root", "admin", "log"}},
	"AdminUsers":              {hierarchy: []string{"root", "admin", "users"}},
	"Categories":              {hierarchy: []string{"root", "_:categories"}},
	"CategoriesForums":        {hierarchy: []string{"root", "_:categories", ":forums"}},
	"CategoriesForumsTopics":  {hierarchy: []string{"root", "_:categories", ":forums", ":topics"}},
	"Docs":                    {hierarchy: []string{"root", "docs"}, noBreadcrumb: true},
	"DocsAbout":               {hierarchy: []string{"root", "docs", "about"}},
	"DocsContact":             {hierarchy: []string{"root", "docs", "contact"}},
	"DocsManual":              {hierarchy: []string{"root", "docs", "manual"}},
	"DocsPrivacy":             {hierarchy: []string{"root", "docs", "privacy"}},
	"DocsTerms":               {hierarchy: []string{"root", "docs", "terms"}},
	"New":                     {hierarchy: []string{"root", "new"}},
	"Search":                  {hierarchy: []string{"root", "search"}},
	"Session":                 {hierarchy: []string{"root", "session"}, noBreadcrumb: true},
	"SessionConfirm":          {hierarchy: []string{"root", "session", "confirm"}},
	"SessionConfirmResend":    {hierarchy: []string{"root", "session", "confirm", "resend"}},
	"SessionSettings":         {hierarchy: []string{"root", "session", "settings"}},
	"SessionSignin":           {hierarchy: []string{"root", "session", "signin"}},
	"SessionSignout":          {hierarchy: []string{"root", "session", "signout"}},
	"SessionSignup":           {hierarchy: []string{"root", "session", "signup"}},
	"SessionProvider":         {hierarchy: []string{"root", "session", ":provider"}},
	"SessionProviderCallback": {hierarchy: []string{"root", "session", ":provider", "callback"}},
	"System":                  {hierarchy: []string{"root", "system"}},
}

func (r Route) Len() int {
	return len(r.hierarchy)
}

func (r Route) AsURL() string {
	if r.Len() <= 1 {
		return ""
	}
	return strings.Join([]string(r.hierarchy)[1:], "/")
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

func (r Route) Fill(params map[string]string) (rt Route) {
	var filledHierarchy []string

	for _, segment := range r.hierarchy {
		cidx := strings.Index(segment, ":")
		if cidx > -1 {
			segmentVar := segment[cidx+1:]
			if filler, ok := params[segmentVar]; ok {
				filledHierarchy = append(filledHierarchy, segment[0:cidx]+filler)
			} else {
				filledHierarchy = append(filledHierarchy, segment)
			}
		} else {
			filledHierarchy = append(filledHierarchy, segment)
		}
	}

	return Route{
		hierarchy: filledHierarchy,
	}
}

func For(id string) (r Route) {
	var ok bool

	if r, ok = Routes[id]; !ok {
		return Route{}
	}

	return r
}
