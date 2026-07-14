package admin

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/session"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) RecordVisit(c fiber.Ctx) (err error) {
	if c.Method() != fiber.MethodGet {
		return c.Next()
	}

	actorID, ok := session.New(c).GetUserUUID()
	if !ok {
		return c.Next()
	}

	role, _ := helpers.CurrentUserRoleAndGroups(r.Runtime, c)
	if role != user.AdminRole {
		return c.Next()
	}

	logicactivity.RecordAdminVisit(r.Runtime, actorID, c.Path())

	return c.Next()
}
