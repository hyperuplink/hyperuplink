package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Vote in a topic's poll
// @Description	A vote is final and can only be cast while the poll is still open,
// @Description	and the results are what the topic hands back afterwards.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			categories	path		string						true	"The category slug"
// @Param			forums		path		string						true	"The forum slug"
// @Param			topics		path		string						true	"The topic slug"
// @Param			request		body		topics.TopicPollVoteBody	true	"The option to vote for"
// @Success		200			{object}	request.StatusResponse
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Failure		422			{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/_{categories}/{forums}/{topics}/poll [post]
func (r *Route) PollVote(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(TopicPollVoteBody)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	err = logictopics.PollVoteBySlugs(r.Runtime, req.Perms(), &logictopics.PollVoteBySlugsInput{
		ForumSlug: c.Params("forums"),
		TopicSlug: c.Params("topics"),
		AuthorID:  authorID,
		Selection: in.Selection,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
