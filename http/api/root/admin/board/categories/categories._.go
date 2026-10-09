package categories

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
	r.Path = route.For("APIV1AdminBoardCategories").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("", r.Create).Name("create")
		base.Put("/"+route.For("APIV1AdminBoardCategory").Pathname(), r.Update).Name("update")
		base.Delete("/"+route.For("APIV1AdminBoardCategory").Pathname(), r.Destroy).Name("destroy")
		base.Post("/"+route.For("APIV1AdminBoardCategory").Pathname()+"/"+route.For("APIV1AdminBoardCategoryMoveUp").Pathname(), r.MoveUp).Name("moveup")
		base.Post("/"+route.For("APIV1AdminBoardCategory").Pathname()+"/"+route.For("APIV1AdminBoardCategoryMoveDown").Pathname(), r.MoveDown).Name("movedown")
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
