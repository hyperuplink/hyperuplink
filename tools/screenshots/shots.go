package main

import (
	"github.com/chromedp/chromedp"
)

type shot struct {
	File  string
	URL   string
	As    string
	Wait  string
	Clip  string
	Below string
	Prep  []chromedp.Action
}

type set struct {
	Name    string
	Out     string
	Width   int
	Height  int
	Clip    string
	Below   string
	Resize  string
	Quality int
	Shots   []shot
}

var sets = []set{manual, site}

var manual = set{
	Name:    "manual",
	Out:     "../../docs/manual",
	Width:   1280,
	Height:  860,
	Clip:    ".container",
	Quality: 55,
	Shots:   manualShots,
}

var site = set{
	Name:    "site",
	Out:     "../../../pub/static/screenshots",
	Width:   1280,
	Height:  800,
	Below:   ".header",
	Resize:  "960x600",
	Quality: 80,
	Shots:   siteShots,
}

var siteShots = []shot{
	{
		File: "forum-overview.webp",
		URL:  "/",
		As:   "vera",
	},
	{
		File: "topics-replies.webp",
		URL:  "/_general/announcements/welcome-to-the-board",
		As:   "vera",
	},
	{
		File: "administration.webp",
		URL:  "/admin",
		As:   "sysop",
	},
	{
		File:  "manual.webp",
		URL:   "/docs/manual",
		As:    "vera",
		Below: ".markdown img",
	},
}

var manualShots = []shot{
	{
		File: "session/signin.webp",
		URL:  "/session/signin",
		As:   "",
	},
	{
		File: "session/signup.webp",
		URL:  "/session/signup",
		As:   "",
	},
	{
		File: "root.webp",
		URL:  "/",
		As:   "vera",
	},
	{
		File: "categories/category.webp",
		URL:  "/_general",
		As:   "vera",
	},
	{
		File: "categories/forum.webp",
		URL:  "/_general/announcements",
		As:   "vera",
	},
	{
		File: "categories/topic.webp",
		URL:  "/_general/announcements/welcome-to-the-board",
		As:   "vera",
	},
	{
		File: "newpost/new-topic.webp",
		URL:  "/new",
		As:   "vera",
	},
	{
		File: "newpost/poll-editor.webp",
		URL:  "/new",
		As:   "vera",
		Prep: []chromedp.Action{
			chromedp.Click(`label[for="kind-poll"]`, chromedp.ByQuery),
		},
	},
	{
		File: "search/search.webp",
		URL:  "/search?query=theme&fields=title&fields=topic&fields=replies",
		As:   "vera",
	},
	{
		File: "report/report.webp",
		URL:  "/report?target=topic&id=rules001",
		As:   "vera",
	},
	{
		File: "user/profile.webp",
		URL:  "/~juno",
		As:   "vera",
	},
	{
		File: "account/account.webp",
		URL:  "/account",
		As:   "vera",
	},
	{
		File: "account/profile/profile.webp",
		URL:  "/account/profile",
		As:   "vera",
	},
	{
		File: "account/password/password.webp",
		URL:  "/account/password",
		As:   "vera",
	},
	{
		File: "account/settings/settings.webp",
		URL:  "/account/settings",
		As:   "vera",
	},
	{
		File: "account/twofactor/twofactor.webp",
		URL:  "/account/twofactor",
		As:   "vera",
	},
	{
		File: "account/api/api.webp",
		URL:  "/account/api",
		As:   "vera",
		Prep: []chromedp.Action{
			chromedp.SendKeys(`input[name="name"]`, "Fediverse bot", chromedp.ByQuery),
			chromedp.Click(`button.affirmative`, chromedp.ByQuery),
			chromedp.WaitVisible(`table.detailed`, chromedp.ByQuery),
		},
	},
	{
		File: "categories/poll-results.webp",
		URL:  "/_support/help/which-theme-should-be-the-default",
		As:   "juno",
	},
	{
		File: "categories/poll.webp",
		URL:  "/_support/help/which-theme-should-be-the-default",
		As:   "sysop",
	},
	{
		File: "admin/admin.webp",
		URL:  "/admin",
		As:   "sysop",
	},
	{
		File: "admin/general/general.webp",
		URL:  "/admin/general",
		As:   "sysop",
	},
	{
		File: "admin/board/board.webp",
		URL:  "/admin/board",
		As:   "sysop",
	},
	{
		File: "admin/board/categories/categories.webp",
		URL:  "/admin/board/categories",
		As:   "sysop",
	},
	{
		File: "admin/board/forums/forums.webp",
		URL:  "/admin/board/forums",
		As:   "sysop",
	},
	{
		File: "admin/board/topics/topics.webp",
		URL:  "/admin/board/topics",
		As:   "sysop",
	},
	{
		File: "admin/board/attachments/attachments.webp",
		URL:  "/admin/board/attachments",
		As:   "sysop",
	},
	{
		File: "admin/board/profiles/profiles.webp",
		URL:  "/admin/board/profiles",
		As:   "sysop",
	},
	{
		File: "admin/board/themes/themes.webp",
		URL:  "/admin/board/themes",
		As:   "sysop",
	},
	{
		File: "admin/permissions/permissions.webp",
		URL:  "/admin/permissions",
		As:   "sysop",
	},
	{
		File: "admin/users/users.webp",
		URL:  "/admin/users",
		As:   "sysop",
	},
	{
		File: "admin/reports/reports.webp",
		URL:  "/admin/reports",
		As:   "sysop",
	},
	{
		File: "admin/auth/auth.webp",
		URL:  "/admin/auth",
		As:   "sysop",
	},
	{
		File: "admin/comms/comms.webp",
		URL:  "/admin/comms",
		As:   "sysop",
	},
	{
		File: "admin/comms/email/email.webp",
		URL:  "/admin/comms/email",
		As:   "sysop",
	},
	{
		File: "admin/comms/xmpp/xmpp.webp",
		URL:  "/admin/comms/xmpp",
		As:   "sysop",
	},
	{
		File: "admin/logs/logs.webp",
		URL:  "/admin/logs",
		As:   "sysop",
	},
	{
		File: "admin/health/health.webp",
		URL:  "/admin/health",
		As:   "sysop",
	},
}
