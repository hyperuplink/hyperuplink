package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show the two-factor state
// @Description	Whether two-factor authentication is on. While it is off, a fresh
// @Description	enrollment comes along, with the otpauth URL to scan, its secret
// @Description	and a QR code as a data URI; the client sends the URL back with a
// @Description	code to turn it on.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{enabled=boolean,enrollment=view.OTPEnrollment}
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/twofactor [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	if req.User().OTPEnabled {
		return req.Respond(fiber.Map{"enabled": true})
	}

	enrollment, err := user.NewOTPEnrollment(req.System.Name, req.User().Username, "")
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"enabled":    false,
		"enrollment": view.OTPEnrollmentFrom(enrollment),
	})
}
