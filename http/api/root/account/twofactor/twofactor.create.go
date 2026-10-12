package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Enable two-factor authentication
// @Description	Turns two-factor authentication on with the otpauth URL of an
// @Description	enrollment and a code generated from it, which proves that the
// @Description	authenticator holds the secret.
// @Tags			account
// @Accept			json
// @Param			request	body	twofactor.EnableBody	true	"The otpauth URL and a current code"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/twofactor [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(EnableBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	secret, err := user.OTPSecretFromURL(in.OTPURL)
	if err != nil {
		return req.RespondProblem(problem.Validation(map[string][]string{"otp_url": {"url"}}))
	}

	err = logictwofactor.Enable(req.Context(), r.Runtime, req.User(), secret, in.OTPCode)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
