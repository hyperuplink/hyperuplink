package profiles

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicprofiles "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/profiles"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Update the profile settings
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		request	body		logicprofiles.UpdateInput	true	"The profile settings to store"
// @Success	200		{object}	request.StatusResponse
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/profiles [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicprofiles.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	actorID := uuid.NullUUID{}
	actorID.UUID, actorID.Valid = req.UserUUID()

	err = logicprofiles.Update(r.Runtime, actorID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
