package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
)

type Route struct {
	route.RouteController
}

type TwofactorEnableInput struct {
	OTPURL  string `json:"otp_url" form:"otp_url" validate:"required"`
	OTPCode string `json:"otp_code" form:"otp_code" validate:"required,numeric,len=6"`
}

type TwofactorDisableInput struct {
	CurrentPassword string `json:"current_password" form:"current_password" validate:"required,min=8,max=64"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("AccountTwofactor").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("show")
		base.Post("/enable", r.Enable).Name("enable")
		base.Post("/disable", r.Disable).Name("disable")
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
