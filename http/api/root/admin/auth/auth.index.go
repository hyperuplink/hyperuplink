package auth

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	settingAuth, err := settingRepo.GetByID[setting.Auth](
		r.Runtime.Repositories.Setting,
		"auth",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	authProviders, err := r.Runtime.Config.AuthProviders()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	providerTypes := []string{}
	for _, provider := range authProviders {
		providerTypes = append(providerTypes, provider.Type)
	}

	return req.Respond(fiber.Map{
		"auth":           settingAuth.JSONValue,
		"auth_providers": providerTypes,
	})
}
