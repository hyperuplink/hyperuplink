package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
)

// @Summary		Set a user's groups
// @Description	Replaces the groups the user belongs to with member_of; ids that
// @Description	name no group are dropped.
// @Tags			admin
// @Accept			json
// @Param			user	path	string						true	"The user UUID"
// @Param			request	body	logicuser.MembershipInput	true	"The group ids"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/users/{user}/groups [put]
func (r *Route) Groups(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("user")
	if ret {
		return rerr
	}

	in := new(logicuser.MembershipInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicuser.UpdateMembership(r.Runtime, id.String(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
