package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountTwofactor")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)

	if !usr.OTPEnabled {
		pendingURL, _ := req.Session.GetPendingOTPURL()

		var enrollment *user.OTPEnrollment
		enrollment, err = user.NewOTPEnrollment(
			req.System.Name,
			usr.Username,
			pendingURL,
		)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}

		req.Session.SetPendingOTPURL(enrollment.URL)
		req.SetData("enrollment", enrollment)
	}

	return req.Respond()
}
