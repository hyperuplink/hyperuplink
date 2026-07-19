package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicroot "xn--gckvb8fzb.com/hyperuplink/logic/root"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("Root")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var board *logicroot.Board
	board, err = logicroot.BoardView(r.Runtime, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("categories_forums", board.CategoriesForums)
	req.SetData("topics", &board.RecentTopics)

	return req.Respond()
}
