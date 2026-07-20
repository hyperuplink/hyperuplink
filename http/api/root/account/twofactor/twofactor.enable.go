package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Enable two-factor authentication
// @Description	The otpauth URL handed out by the enrollment is sent back together
// @Description	with a code generated from it, which proves the authenticator app
// @Description	holds the same secret.
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			request	body		twofactor.TwofactorEnableInput	true	"The otpauth URL and a code generated from it"
// @Success		200		{object}	request.StatusResponse
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/twofactor/enable [post]
func (r *Route) Enable(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(TwofactorEnableInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	var secret string
	secret, err = user.OTPSecretFromURL(in.OTPURL)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = logictwofactor.Enable(r.Runtime, req.User, secret, in.OTPCode)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
