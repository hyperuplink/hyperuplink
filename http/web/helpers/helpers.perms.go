package helpers

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/session"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func CurrentUserRoleAndGroups(
	rt *runtime.Runtime,
	c fiber.Ctx,
) (role user.Role, memberOf []string) {
	s := session.New(c)

	userID, ok := s.GetUserID()
	if !ok {
		return user.GuestRole, nil
	}

	usr, err := gh.Repositories(rt).User.GetByID(userID, common.QueryOptions{
		WithBanned:  false,
		WithSpammed: false,
		WithDeleted: false,
		Limit:       1,
	})
	if err != nil {
		return user.GuestRole, nil
	}

	return usr.Role, usr.MemberOf
}
