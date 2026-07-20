package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/attachments"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/categories"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/docs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/report"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/search"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/session"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/topics"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/user"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
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
	r.Path = ""
	r.Env = route.NewEnv()

	r.Router.Route("/", func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		sessionRoute, err := session.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionRoute)

		accountRoute, err := account.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountRoute)

		adminRoute, err := admin.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminRoute)

		attachmentsRoute, err := attachments.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, attachmentsRoute)

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)

		topicsRoute, err := topics.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, topicsRoute)

		docsRoute, err := docs.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, docsRoute)

		newpostRoute, err := newpost.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, newpostRoute)

		searchRoute, err := search.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, searchRoute)

		reportRoute, err := report.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, reportRoute)

		userRoute, err := user.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, userRoute)
	}, "root.")

	return r, nil
}

func (r *Route) GetRuntime() *runtime.Runtime {
	return r.Runtime
}

func (r *Route) GetPath() string {
	return r.Path
}

func (r *Route) GetEnv() *route.Environment {
	return r.Env
}
