package account

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
)

// @Summary		Update the account
// @Description	Changes the fields that are present: the timezone, the signature,
// @Description	which is Markdown, and whether replies to the user's topics are
// @Description	announced by message. Fields that are left out stay as they are.
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			request	body		account.PatchBody	true	"The fields to change"
// @Success		200		{object}	object{user=view.Account}
// @Failure		401		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account [patch]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(PatchBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	usr := req.User()

	acc, err := logicaccount.Show(r.Runtime, usr)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	if in.Timezone != nil {
		err = logicsettings.Update(r.Runtime, usr, &logicsettings.UpdateInput{Timezone: *in.Timezone})
		if ret, rerr := req.RespondOnError(err); ret {
			return rerr
		}
	}

	if in.SignatureText != nil || in.NotifyOnReply != nil {
		profile := &logicprofile.UpdateInput{
			SignatureText: usr.SignatureText,
			NotifyOnReply: acc.Profile.NotifyOnReply,
		}
		if in.SignatureText != nil {
			profile.SignatureText = *in.SignatureText
		}
		if in.NotifyOnReply != nil {
			profile.NotifyOnReply = *in.NotifyOnReply
		}

		err = logicprofile.Update(r.Runtime, usr, profile)
		if ret, rerr := req.RespondOnError(err); ret {
			return rerr
		}
	}

	return r.respondAccount(req)
}
