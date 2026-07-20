package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show the two-factor state
// @Description	When two-factor authentication is off this hands out a fresh
// @Description	enrollment with the secret and the otpauth URL to scan, and when it
// @Description	is already on only the enabled flag comes back.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{enabled=boolean,enrollment=user.OTPEnrollment}
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/twofactor [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	if req.User.OTPEnabled {
		return req.Respond(fiber.Map{
			"enabled": true,
		})
	}

	var enrollment *user.OTPEnrollment
	enrollment, err = user.NewOTPEnrollment(
		req.System.Name,
		req.User.Username,
		"",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"enabled":    false,
		"enrollment": enrollment,
	})
}
