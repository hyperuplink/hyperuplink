package board

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/attachments"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/categories"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/profiles"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/themes"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board/topics"
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
	r.Path = route.For("AdminBoard").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		attachmentsRoute, err := attachments.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, attachmentsRoute)

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)

		forumsRoute, err := forums.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, forumsRoute)

		profilesRoute, err := profiles.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, profilesRoute)

		themesRoute, err := themes.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, themesRoute)

		topicsRoute, err := topics.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, topicsRoute)
	}, r.Path+".")

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
