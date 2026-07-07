package user

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("User")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "user/index",
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Maybe remove!
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var username string = c.Params("user")
	req.UpdateTitle("~" + username)

	usr, err := r.Runtime.Repositories.User.GetByUsername(
		username,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)

	groups, err := r.Runtime.Repositories.Group.All(common.QueryOptions{
		OrderBy: "name",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("groups", groups)

	topics, err := r.Runtime.Repositories.Topic.VAllForAuthorUUID(
		usr.ID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
			Limit:   10,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("topics", topics)

	replies, err := r.Runtime.Repositories.Reply.VAllWithTopicForAuthorUUID(
		usr.ID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
			Limit:   10,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("replies", replies)

	return req.Respond()
}
