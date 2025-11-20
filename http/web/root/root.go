package root

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root/categories"
	"github.com/mrusme/hyperuplink/http/web/site"
	"github.com/mrusme/hyperuplink/runtime"
)

type Route struct {
	route.Route
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router

	r.Router.Route("", func(base fiber.Router) {
		base.Get("/", r.Index).Name("index")

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)
	}, "root.")

	return r, nil
}

func (r *Route) Index(c fiber.Ctx) error {
	err := c.App().ReloadViews()
	r.Runtime.Error("error", err)
	return c.Render("views/root", fiber.Map{
		"Site": site.New(r.Runtime, c, "Root"),
	}, "views/layouts/base")
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
