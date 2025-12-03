package request

import (
	"slices"

	"github.com/mrusme/hyperuplink/models/user"
)

func (req *Request) AccessControl(roles ...user.Role) (mustReturn bool, err error) {
	if slices.Index(roles, req.Session.GetCurrentUserRole()) < 0 {
		return true, req.RedirectToRoot()
	}

	return false, nil
}
