package categories

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/categories"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show a category
// @Description	A category the caller may not read is answered as though it did not
// @Description	exist, so that the board gives nothing away about what it hides.
// @Tags			board
// @Produce		json
// @Param			categories	path		string	true	"The category slug"
// @Success		200			{object}	logiccategories.View
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/_{categories} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logiccategories.View
	view, err = logiccategories.Show(r.Runtime, c.Params("categories"), req.Perms())
	if errors.Is(err, errs.ErrForbidden) {
		return req.RespondError(errs.ErrNoRows)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(view)
}
