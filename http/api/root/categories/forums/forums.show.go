package forums

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show a forum
// @Description	Lists one page of the forum's topics, newest activity first.
// @Tags			board
// @Produce		json
// @Param			categories	path		string	true	"The category slug"
// @Param			forums		path		string	true	"The forum slug"
// @Param			page		query		integer	false	"The page to return, counting from one"
// @Success		200			{object}	logicforums.View
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Failure		422			{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/_{categories}/{forums} [get]
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

	var view *logicforums.View
	view, err = logicforums.Show(r.Runtime, &logicforums.ShowInput{
		ForumSlug: c.Params("forums"),
		Page:      activePage,
		PerPage:   req.System.GetTopicsPerPage(),
		ViewerID:  viewerID,
	}, req.Perms())
	if errors.Is(err, errs.ErrForbidden) {
		return req.RespondError(errs.ErrNoRows)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(view)
}
