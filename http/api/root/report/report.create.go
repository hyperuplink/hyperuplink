package report

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicreport "xn--gckvb8fzb.com/hyperuplink/logic/root/report"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Report a post
// @Description	The report is filed for the administrators to read under Admin ->
// @Description	Reports.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			request	body		logicreport.CreateInput	true	"The post and the reason to report it"
// @Success		201		{object}	request.StatusResponse
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		404		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/report [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicreport.CreateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	err = logicreport.Create(r.Runtime, authorID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondCreated(fiber.Map{"status": "ok"})
}
