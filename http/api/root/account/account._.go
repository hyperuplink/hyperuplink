package account

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account/keys"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account/password"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account/picture"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/account/twofactor"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
)

type Route struct {
	route.RouteController
}

type PatchBody struct {
	Timezone      *string `json:"timezone" validate:"omitempty,timezone"`
	SignatureText *string `json:"signature_text" validate:"omitempty,max=256"`
	NotifyOnReply *bool   `json:"notify_on_reply"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("APIV1Account").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Patch("", r.Update).Name("update")

		keysRoute, err := keys.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, keysRoute)

		passwordRoute, err := password.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, passwordRoute)

		pictureRoute, err := picture.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, pictureRoute)

		twofactorRoute, err := twofactor.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, twofactorRoute)
	}, r.Path+".")

	return r, nil
}

func (r *Route) respondAccount(req *request.Request) error {
	usr, err := gh.Repositories(r.Runtime).User.GetByUUID(req.User().ID, common.QueryOptions{Limit: 1})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	acc, err := logicaccount.Show(r.Runtime, usr)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"user": req.Views().Account(usr, acc.Groups, acc.Profile)})
}
