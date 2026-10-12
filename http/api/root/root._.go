package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/attachments"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/board"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/categories"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/forums"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/pages"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/replies"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/reports"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/search"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/session"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/topics"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/users"
)

type Route struct {
	route.RouteController
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("APIV1").AsURL()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		sessionRoute, err := session.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionRoute)

		boardRoute, err := board.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, boardRoute)

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)

		forumsRoute, err := forums.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, forumsRoute)

		topicsRoute, err := topics.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, topicsRoute)

		repliesRoute, err := replies.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, repliesRoute)

		reportsRoute, err := reports.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, reportsRoute)

		searchRoute, err := search.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, searchRoute)

		usersRoute, err := users.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, usersRoute)

		pagesRoute, err := pages.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, pagesRoute)

		attachmentsRoute, err := attachments.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, attachmentsRoute)

		accountRoute, err := account.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountRoute)

		adminRoute, err := admin.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminRoute)
	}, "api.v1.")

	return r, nil
}
