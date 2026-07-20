package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Create a forum
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		request	body		logicforums.CreateInput	true	"The forum to create"
// @Success	201		{object}	object{id=string}
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/forums [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicforums.CreateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	id, err := logicforums.Create(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondCreated(fiber.Map{"id": id})
}
