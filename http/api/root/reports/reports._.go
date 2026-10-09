package reports

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
)

type Route struct {
	route.RouteController
}

type ReportBody struct {
	Target string `json:"target" validate:"required,oneof=topic reply"`
	ID     string `json:"id" validate:"required,max=64"`
	Reason string `json:"reason" validate:"required,oneof=spam misconduct illegal"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("APIV1Reports").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Post("", r.Create).Name("create")
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
