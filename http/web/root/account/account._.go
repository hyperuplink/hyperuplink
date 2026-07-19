package account

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account/api"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account/password"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account/profile"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account/twofactor"
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
	r.Path = route.For("Account").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		accountProfileRoute, err := profile.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountProfileRoute)

		accountSettingsRoute, err := settings.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountSettingsRoute)

		accountPasswordRoute, err := password.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountPasswordRoute)

		accountTwofactorRoute, err := twofactor.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountTwofactorRoute)

		accountAPIRoute, err := api.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountAPIRoute)
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
