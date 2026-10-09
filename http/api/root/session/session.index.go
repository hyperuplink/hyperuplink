package session

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
)

// @Summary		Show the viewer
// @Description	Who the request acts as: a guest when no key was sent, otherwise
// @Description	the key's user with the key itself, and in both cases the
// @Description	permission level the viewer has by default and per category.
// @Tags			session
// @Produce		json
// @Success		200	{object}	view.Session
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/session [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	sess := view.Session{
		Authenticated: !req.IsGuest(),
		Role:          string(req.Role()),
		Permissions:   view.SessionPermissionsFor(req.Perms()),
	}

	if usr := req.User(); usr != nil {
		var acc *logicaccount.View
		acc, err = logicaccount.Show(r.Runtime, usr)
		if ret, rerr := req.RespondOnError(err); ret {
			return rerr
		}

		account := req.Views().Account(usr, acc.Groups, acc.Profile)
		sess.User = &account

		key := view.KeyFrom(req.Principal.Key)
		sess.Key = &key
	}

	return req.Respond(sess)
}
