package root

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/http/web/root/account"
	"github.com/mrusme/hyperuplink/http/web/root/categories"
	"github.com/mrusme/hyperuplink/http/web/root/sessions"
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
	r.Path = ""
	r.Env = route.NewEnv()

	r.Router.Route("/", func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		sessionsRoute, err := sessions.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionsRoute)

		accountRoute, err := account.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountRoute)

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)
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

func (r *Route) Index(c fiber.Ctx) error {
	req := request.New(r, c, []string{"base"}, "root", "root_noun")

	err := c.App().ReloadViews()
	r.Runtime.Error("error", err)
	return req.Respond()
}

func (r *Route) Show(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Destroy(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}
