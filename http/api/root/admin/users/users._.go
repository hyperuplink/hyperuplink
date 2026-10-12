package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicusers "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/users"
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
	r.Path = route.For("APIV1AdminUsers").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("/"+route.For("APIV1AdminUser").Pathname()+"/"+route.For("APIV1AdminUserConfirm").Pathname(), r.Confirm).Name("confirm")
		base.Post("/"+route.For("APIV1AdminUser").Pathname()+"/"+route.For("APIV1AdminUserBan").Pathname(), r.Ban).Name("ban")
		base.Post("/"+route.For("APIV1AdminUser").Pathname()+"/"+route.For("APIV1AdminUserUnban").Pathname(), r.Unban).Name("unban")
		base.Put("/"+route.For("APIV1AdminUser").Pathname()+"/"+route.For("APIV1AdminUserGroups").Pathname(), r.Groups).Name("groups")
		base.Delete("/"+route.For("APIV1AdminUser").Pathname(), r.Destroy).Name("destroy")
	}, r.Path+".")

	return r, nil
}

func (r *Route) action(
	c fiber.Ctx,
	fn func(*runtime.Runtime, *logicusers.UserInput) error,
) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("user")
	if ret {
		return rerr
	}

	err = fn(r.Runtime, &logicusers.UserInput{ID: id.String()})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
