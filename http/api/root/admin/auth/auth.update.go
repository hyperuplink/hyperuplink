package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicauth "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/auth"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Update the authentication settings
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		request	body		logicauth.UpdateInput	true	"The authentication settings to store"
// @Success	200		{object}	request.StatusResponse
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/auth [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicauth.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	actorID := uuid.NullUUID{}
	actorID.UUID, actorID.Valid = req.UserUUID()

	err = logicauth.Update(r.Runtime, actorID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
