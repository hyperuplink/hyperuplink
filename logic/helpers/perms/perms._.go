package perms

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func Resolve(
	rt *runtime.Runtime,
	role user.Role,
	memberOf []string,
) *permission.Resolution {
	if role == user.AdminRole {
		return &permission.Resolution{IsAdmin: true}
	}

	cats, err := gh.Repositories(rt).Category.All(common.QueryOptions{})
	if err != nil {
		rt.Error("perms", err)
		return permission.NewResolution()
	}

	res, err := gh.Repositories(rt).Permission.Resolve(role, memberOf, cats)
	if err != nil {
		rt.Error("perms", err)
		return permission.NewResolution()
	}

	return res
}
