package routes

import (
	"xn--gckvb8fzb.com/glides/http/route"
)

const DEFAULT_TITLE string = "Hyperuplink"

// TODO: Find a way to build this automatically
var Routes route.Table = route.Table{
	"Root":                    {Hierarchy: []string{"root"}},
	"Account":                 {Hierarchy: []string{"root", "account"}},
	"AccountAPI":              {Hierarchy: []string{"root", "account", "api"}},
	"AccountPassword":         {Hierarchy: []string{"root", "account", "password"}},
	"AccountProfile":          {Hierarchy: []string{"root", "account", "profile"}},
	"AccountSettings":         {Hierarchy: []string{"root", "account", "settings"}},
	"AccountSettingsView":     {Hierarchy: []string{"root", "account", "settings", "view"}, NoBreadcrumb: true},
	"AccountTwofactor":        {Hierarchy: []string{"root", "account", "twofactor"}},
	"Admin":                   {Hierarchy: []string{"root", "admin"}},
	"Attachment":              {Hierarchy: []string{"root", "attachment", ":attachment"}},
	"Attachments":             {Hierarchy: []string{"root", "attachments"}, NoBreadcrumb: true},
	"AdminAuth":               {Hierarchy: []string{"root", "admin", "auth"}},
	"AdminBoard":              {Hierarchy: []string{"root", "admin", "board"}},
	"AdminBoardAttachments":   {Hierarchy: []string{"root", "admin", "board", "attachments"}},
	"AdminBoardCategories":    {Hierarchy: []string{"root", "admin", "board", "categories"}},
	"AdminBoardForums":        {Hierarchy: []string{"root", "admin", "board", "forums"}},
	"AdminBoardTopics":        {Hierarchy: []string{"root", "admin", "board", "topics"}},
	"AdminBoardProfiles":      {Hierarchy: []string{"root", "admin", "board", "profiles"}},
	"AdminBoardThemes":        {Hierarchy: []string{"root", "admin", "board", "themes"}},
	"AdminComms":              {Hierarchy: []string{"root", "admin", "comms"}},
	"AdminCommsEmail":         {Hierarchy: []string{"root", "admin", "comms", "email"}},
	"AdminCommsXmpp":          {Hierarchy: []string{"root", "admin", "comms", "xmpp"}},
	"AdminGeneral":            {Hierarchy: []string{"root", "admin", "general"}},
	"AdminHealth":             {Hierarchy: []string{"root", "admin", "health"}},
	"AdminLogs":               {Hierarchy: []string{"root", "admin", "logs"}},
	"AdminPermissions":        {Hierarchy: []string{"root", "admin", "permissions"}},
	"AdminReports":            {Hierarchy: []string{"root", "admin", "reports"}},
	"AdminUsers":              {Hierarchy: []string{"root", "admin", "users"}},
	"Categories":              {Hierarchy: []string{"root", "_:categories"}},
	"CategoriesForums":        {Hierarchy: []string{"root", "_:categories", ":forums"}},
	"CategoriesForumsTopics":  {Hierarchy: []string{"root", "_:categories", ":forums", ":topics"}},
	"Docs":                    {Hierarchy: []string{"root", "docs"}, NoBreadcrumb: true},
	"DocsAbout":               {Hierarchy: []string{"root", "docs", "about"}},
	"DocsContact":             {Hierarchy: []string{"root", "docs", "contact"}},
	"DocsManual":              {Hierarchy: []string{"root", "docs", "manual"}},
	"DocsPrivacy":             {Hierarchy: []string{"root", "docs", "privacy"}},
	"DocsTerms":               {Hierarchy: []string{"root", "docs", "terms"}},
	"New":                     {Hierarchy: []string{"root", "new"}},
	"Report":                  {Hierarchy: []string{"root", "report"}},
	"Search":                  {Hierarchy: []string{"root", "search"}},
	"Session":                 {Hierarchy: []string{"root", "session"}, NoBreadcrumb: true},
	"SessionConfirm":          {Hierarchy: []string{"root", "session", "confirm"}},
	"SessionConfirmResend":    {Hierarchy: []string{"root", "session", "confirm", "resend"}},
	"SessionSignin":           {Hierarchy: []string{"root", "session", "signin"}},
	"SessionSignout":          {Hierarchy: []string{"root", "session", "signout"}},
	"SessionSignup":           {Hierarchy: []string{"root", "session", "signup"}},
	"SessionTwofactor":        {Hierarchy: []string{"root", "session", "twofactor"}},
	"SessionProvider":         {Hierarchy: []string{"root", "session", ":provider"}},
	"SessionProviderCallback": {Hierarchy: []string{"root", "session", ":provider", "callback"}},
	"Topics":                  {Hierarchy: []string{"root", "topics"}, NoBreadcrumb: true},
	"User":                    {Hierarchy: []string{"root", "~:user"}},
	"System":                  {Hierarchy: []string{"root", "system"}},
}

func Use() {
	route.Use(Routes, DEFAULT_TITLE)
}
