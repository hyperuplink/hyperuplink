package helpers

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/session"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
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

	usr, err := rt.Repositories.User.GetByID(userID, common.QueryOptions{
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

func ResolvePermissions(
	rt *runtime.Runtime,
	role user.Role,
	memberOf []string,
) *permission.Resolution {
	if role == user.AdminRole {
		return &permission.Resolution{IsAdmin: true}
	}

	cats, err := rt.Repositories.Category.All(common.QueryOptions{})
	if err != nil {
		rt.Error("perms", err)
		return permission.NewResolution()
	}

	res, err := rt.Repositories.Permission.Resolve(role, memberOf, cats)
	if err != nil {
		rt.Error("perms", err)
		return permission.NewResolution()
	}

	return res
}
