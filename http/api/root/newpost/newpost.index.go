package newpost

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show what a new post may be written into
// @Description	Hands back the forums the caller may post in, and, when a forum,
// @Description	topic or reply is named in the query, the post that is being
// @Description	answered.
// @Tags			board
// @Produce		json
// @Param			forum	query		string	false	"The slug of the forum to post in"
// @Param			topic	query		string	false	"The slug of the topic to reply to"
// @Param			reply	query		string	false	"The identifier of the reply to quote"
// @Success		200		{object}	logicnewpost.FormView
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		404		{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/new [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicnewpost.FormView
	view, err = logicnewpost.View(r.Runtime, req.Perms(), &logicnewpost.FormViewInput{
		ForumSlug: c.Query("forum"),
		TopicSlug: c.Query("topic"),
		ReplyID:   c.Query("reply"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(view)
}
