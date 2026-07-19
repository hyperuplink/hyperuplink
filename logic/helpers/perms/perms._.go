package perms

import (
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func Resolve(
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
