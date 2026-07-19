package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	users, err := r.Runtime.Repositories.User.All(common.QueryOptions{
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
