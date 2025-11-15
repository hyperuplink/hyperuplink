package categories

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
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

	r.Router.Route("/", func(base fiber.Router) {
		base.Get(":slug", r.Show).Name("show")
	}, "categories.")

	return r, nil
}

func (r *Route) Index(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Show(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Destroy(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}
