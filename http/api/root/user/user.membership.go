package user

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Set a user's group membership
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		user	path		string						true	"The username"
// @Param		request	body		logicuser.MembershipInput	true	"The groups the user belongs to"
// @Success	200		{object}	request.StatusResponse
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	404		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/~{user}/membership [post]
func (r *Route) Membership(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicuser.MembershipInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicuser.UpdateMembership(r.Runtime, c.Params("user"), in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
