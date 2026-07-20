package topics

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		List topics
// @Description	The board-wide feed of topics, newest activity first, narrowed to one
// @Description	forum when forum_id is given.
// @Tags			board
// @Produce		json
// @Param			forum_id	query		string	false	"Only list topics in this forum"
// @Param			page		query		integer	false	"The page to return, counting from one"
// @Success		200			{object}	topicsfeed.View
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Failure		422			{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/topics [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var activePage int
	if activePage, err = strconv.Atoi(c.Query("page", "1")); err != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	forumID := uuid.NullUUID{}
	if raw := c.Query("forum_id"); raw != "" {
		id, perr := uuid.Parse(raw)
		if perr != nil {
			return req.RespondError(errs.ErrFormInvalid)
		}
		forumID.UUID, forumID.Valid = id, true
	}

	viewerID := uuid.NullUUID{}
	viewerID.UUID, viewerID.Valid = req.UserUUID()

	var view *topicsfeed.View
	view, err = topicsfeed.Index(r.Runtime, &topicsfeed.Input{
		ForumID:  forumID,
		Page:     activePage,
		PerPage:  req.System.GetTopicsPerPage(),
		ViewerID: viewerID,
	}, req.Perms())
	if errors.Is(err, errs.ErrForbidden) {
		return req.RespondError(errs.ErrNoRows)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(view)
}
