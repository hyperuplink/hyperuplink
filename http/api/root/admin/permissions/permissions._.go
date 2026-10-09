package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
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
	r.Path = route.For("APIV1AdminPermissions").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("/"+route.For("APIV1AdminPermissionsApply").Pathname(), r.Apply).Name("apply")
		base.Post("/"+route.For("APIV1AdminPermissionsRemove").Pathname(), r.Remove).Name("remove")
		base.Post("/"+route.For("APIV1AdminPermissionsGroups").Pathname(), r.GroupCreate).Name("groups.create")
		base.Delete("/"+route.For("APIV1AdminPermissionsGroups").Pathname()+"/"+route.For("APIV1AdminPermissionsGroup").Pathname(), r.GroupDestroy).Name("groups.destroy")
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
