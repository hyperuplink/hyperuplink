package account

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
)

// @Summary		Update the account
// @Description	Changes the fields that are present: the language, one of de, en,
// @Description	es, fr, it and ro, the timezone, the signature, which is Markdown,
// @Description	and whether replies to the user's topics are announced by message.
// @Description	Fields that are left out stay as they are.
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

	acc, err := logicaccount.Show(req.Context(), r.Runtime, usr)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	if in.Language != nil || in.Timezone != nil {
		settings := &logicsettings.UpdateInput{
			Language: usr.Language,
			Timezone: usr.Timezone,
		}
		if in.Language != nil {
			settings.Language = *in.Language
		}
		if in.Timezone != nil {
			settings.Timezone = *in.Timezone
		}
		err = logicsettings.Update(req.Context(), r.Runtime, usr, settings)
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

		err = logicprofile.Update(req.Context(), r.Runtime, usr, profile)
		if ret, rerr := req.RespondOnError(err); ret {
			return rerr
		}
	}

	return r.respondAccount(req)
}
