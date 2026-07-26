package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		List the users
// @Description	Banned and deleted users are listed as well, which is what separates
// @Description	this from the public user endpoint.
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{users=[]user.Detail}
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/admin/users [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	users, err := gh.Repositories(r.Runtime).User.All(common.QueryOptions{
		WithBanned:  true,
		WithDeleted: true,
		OrderBy:     "created_at",
		Order:       common.Descending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	details := make([]user.Detail, 0, len(*users))
	for i := range *users {
		details = append(details, (*users)[i].AsDetail())
	}

	return req.Respond(fiber.Map{
		"users": details,
	})
}
