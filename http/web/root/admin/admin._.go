package admin

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root/admin/auth"
	"github.com/mrusme/hyperuplink/http/web/root/admin/board"
	"github.com/mrusme/hyperuplink/http/web/root/admin/comms"
	"github.com/mrusme/hyperuplink/http/web/root/admin/general"
	"github.com/mrusme/hyperuplink/http/web/root/admin/logs"
	"github.com/mrusme/hyperuplink/http/web/root/admin/users"
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
	r.Path = route.For("Admin").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		adminAuthRoute, err := auth.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminAuthRoute)

		adminBoardRoute, err := board.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminBoardRoute)

		adminCommsRoute, err := comms.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminCommsRoute)

		adminGeneralRoute, err := general.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminGeneralRoute)

		adminLogsRoute, err := logs.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminLogsRoute)

		adminUsersRoute, err := users.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminUsersRoute)
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
