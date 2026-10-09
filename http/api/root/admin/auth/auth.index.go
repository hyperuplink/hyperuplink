package auth

import (
	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

// @Summary		Show the authentication settings
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/auth [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	settingAuth, err := settingRepo.GetByID[setting.Auth](gh.Repositories(r.Runtime).Setting, "auth")
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	var authProviders gh.AuthProviders
	_, err = r.Runtime.Config().Unmarshal("AuthProvider", &authProviders)
	if ret, rerr := req.RespondOnError(err); ret {
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
