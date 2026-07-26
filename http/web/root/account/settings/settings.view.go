package settings

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) View(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountSettingsView")
	parentRoute := route.For("AccountSettings")
	req := request.New(r, c, myRoute,
		[]string{"base"}, parentRoute.AsURL()+"/index",
		parentRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RedirectBackOnError(err); ret == true {
		return rerr
	}

	if _, err = logicsettings.ToggleView(
		r.Runtime,
		usr,
		c.Params("view"),
	); err != nil {
		if errors.Is(err, errs.ErrValidation) {
			return req.RedirectBack()
		}
		if ret, rerr := req.RedirectBackOnError(err); ret == true {
			return rerr
		}
	}

	return req.RedirectBack()
}
