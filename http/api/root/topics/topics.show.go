package topics

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show a topic by identifier
// @Description	The same view the slug path returns, addressed by identifier instead,
// @Description	which saves a caller that already holds one from knowing the slugs.
// @Tags			board
// @Produce		json
// @Param			id		path		string	true	"The topic identifier"
// @Param			page	query		integer	false	"The page to return, counting from one"
// @Success		200		{object}	logictopics.View
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		404		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/topics/{id} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
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

	viewerID := uuid.NullUUID{}
	viewerID.UUID, viewerID.Valid = req.UserUUID()

	var view *logictopics.View
	view, err = logictopics.ShowByID(r.Runtime, &logictopics.ShowByIDInput{
		ID:       c.Params("id"),
		Page:     activePage,
		PerPage:  req.System.GetPostsPerPage(),
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
