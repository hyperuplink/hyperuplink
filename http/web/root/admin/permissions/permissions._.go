package permissions

import (
	"github.com/gofiber/fiber/v3"
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
	r.Path = route.For("AdminPermissions").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("/group", r.GroupCreate).Name("group_create")
		base.Post("/group/destroy", r.GroupDestroy).Name("group_destroy")
		base.Post("/apply", r.Apply).Name("apply")
		base.Post("/remove", r.PermissionRemove).Name("remove")
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
