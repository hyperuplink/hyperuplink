package board

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/attachments"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/categories"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/forums"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/profiles"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/themes"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board/topics"
	"github.com/mrusme/hyperuplink/runtime"
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

		adminBoardAttachmentsRoute, err := attachments.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardAttachmentsRoute)

		adminBoardCategoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardCategoriesRoute)

		adminBoardForumsRoute, err := forums.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardForumsRoute)

		adminBoardProfilesRoute, err := profiles.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardProfilesRoute)

		adminBoardThemesRoute, err := themes.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardThemesRoute)

		adminBoardTopicsRoute, err := topics.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardTopicsRoute)
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
