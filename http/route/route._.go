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
	"Attachment":              {hierarchy: []string{"root", "attachment", ":attachment"}},
	"AdminAuth":               {hierarchy: []string{"root", "admin", "auth"}},
	"AdminBoard":              {hierarchy: []string{"root", "admin", "board"}},
	"AdminBoardAttachments":   {hierarchy: []string{"root", "admin", "board", "attachments"}},
	"AdminBoardCategories":    {hierarchy: []string{"root", "admin", "board", "categories"}},
	"AdminBoardForums":        {hierarchy: []string{"root", "admin", "board", "forums"}},
	"AdminBoardTopics":        {hierarchy: []string{"root", "admin", "board", "topics"}},
	"AdminBoardProfiles":      {hierarchy: []string{"root", "admin", "board", "profiles"}},
	"AdminBoardThemes":        {hierarchy: []string{"root", "admin", "board", "themes"}},
	"AdminComms":              {hierarchy: []string{"root", "admin", "comms"}},
	"AdminCommsEmail":         {hierarchy: []string{"root", "admin", "comms", "email"}},
	"AdminCommsXmpp":          {hierarchy: []string{"root", "admin", "comms", "xmpp"}},
	"AdminGeneral":            {hierarchy: []string{"root", "admin", "general"}},
	"AdminHealth":             {hierarchy: []string{"root", "admin", "health"}},
	"AdminLogs":               {hierarchy: []string{"root", "admin", "logs"}},
	"AdminPermissions":        {hierarchy: []string{"root", "admin", "permissions"}},
	"AdminReports":            {hierarchy: []string{"root", "admin", "reports"}},
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
	"Report":                  {hierarchy: []string{"root", "report"}},
	"Search":                  {hierarchy: []string{"root", "search"}},
	"Session":                 {hierarchy: []string{"root", "session"}, noBreadcrumb: true},
	"SessionConfirm":          {hierarchy: []string{"root", "session", "confirm"}},
	"SessionConfirmResend":    {hierarchy: []string{"root", "session", "confirm", "resend"}},
	"SessionSignin":           {hierarchy: []string{"root", "session", "signin"}},
	"SessionSignout":          {hierarchy: []string{"root", "session", "signout"}},
	"SessionSignup":           {hierarchy: []string{"root", "session", "signup"}},
	"SessionTwofactor":        {hierarchy: []string{"root", "session", "twofactor"}},
	"SessionProvider":         {hierarchy: []string{"root", "session", ":provider"}},
	"SessionProviderCallback": {hierarchy: []string{"root", "session", ":provider", "callback"}},
	"User":                    {hierarchy: []string{"root", "~:user"}},
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

func (r Route) AsID() string {
	joined := strings.Join([]string(r.hierarchy), "_")
	if strings.Index(joined, ":") > -1 {
		return ""
	}
	return joined
}

func (r Route) AsTitle() (title string) {
	if title = r.AsID(); title == "" {
		return title
	}
	return title + "_title"
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

func (r *Route) SetHierarchy(h []string) {
	r.hierarchy = h
}

func For(id string) (r Route) {
	var ok bool

	if r, ok = Routes[id]; !ok {
		return Route{}
	}

	return r
}

func CollidesWithRoute(path string) bool {
	seg := strings.TrimPrefix(path, "/")
	if idx := strings.Index(seg, "/"); idx > -1 {
		seg = seg[:idx]
	}
	seg = strings.ToLower(seg)
	if seg == "" {
		return true
	}

	for _, r := range Routes {
		if r.Len() < 2 {
			continue
		}

		top := strings.ToLower(r.hierarchy[1])
		cidx := strings.Index(top, ":")
		if cidx == -1 {
			if seg == top {
				return true
			}
			continue
		}

		prefix := top[:cidx]
		if prefix == "" || strings.HasPrefix(seg, prefix) {
			return true
		}
	}

	return false
}
