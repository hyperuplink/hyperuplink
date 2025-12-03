package request

import (
	"slices"
)

const (
	GuestRole string = "guest"
	UserRole  string = "user"
	ModRole   string = "mod"
	AdminRole string = "admin"
)

func (req *Request) AccessControl(roles ...string) (mustReturn bool, err error) {
	if slices.Index(roles, req.Site.GetCurrentUserRole()) < 0 {
		return true, req.RedirectToRoot()
	}

	return false, nil
}
