package provider

import (
	"github.com/gofiber/fiber/v3"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
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
	r.Path = route.For("SessionProvider").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("/",
			goth_fiber.BeginAuthHandler).Name("provider.show")
		base.Get("/callback",
			r.ProviderCallbackShow).Name("provider.callback.show")
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

func (r *Route) ProviderCallbackShow(c fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c,
		goth_fiber.CompleteUserAuthOptions{
			ShouldLogout: false,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(user.Email)
}
