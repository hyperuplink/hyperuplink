package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

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
