package user

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show a user
// @Description	The public view of a user, with the groups they belong to and the
// @Description	posts of theirs the caller is allowed to read.
// @Tags			board
// @Produce		json
// @Param			user	path		string	true	"The username"
// @Success		200		{object}	object{user=user.Public,groups=[]group.Group,topics=[]vtopic.VTopic,replies=[]vreply.VReplyWithTopic}
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		404		{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/~{user} [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicuser.View
	view, err = logicuser.Show(r.Runtime, c.Params("user"), req.Perms())
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"user":    view.User.AsPublic(),
		"groups":  view.Groups,
		"topics":  view.Topics,
		"replies": view.Replies,
	})
}
