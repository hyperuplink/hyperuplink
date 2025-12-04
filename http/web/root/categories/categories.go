package categories

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root/categories/forums"
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
	r.Path = route.For("Categories").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/~:slug", func(base fiber.Router) {
		base.Get("/", r.Show).Name("show")

		categoriesForumsRoute, err := forums.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesForumsRoute)
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

func (r *Route) Index(c fiber.Ctx) error {
	return c.SendString("I'm an INDEX request!")
}

func (r *Route) Show(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendString("I'm a CREATE request!")
}

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendString("I'm a UPDATE request!")
}

func (r *Route) Destroy(c fiber.Ctx) error {
	return c.SendString("I'm a DESTROY request!")
}
