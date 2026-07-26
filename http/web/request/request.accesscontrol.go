package request

import (
	"slices"

	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (req *Request) AccessControl(roles ...user.Role) (mustReturn bool, err error) {
	if slices.Index(roles, req.Session.GetCurrentUserRole()) < 0 {
		return true, req.RedirectToRoot()
	}

	return false, nil
}

func (req *Request) GetUser() (usr *user.User, err error) {
	var userID string
	var ok bool
	if userID, ok = req.Session.GetUserID(); !ok {
		return nil, errs.ErrUserIDNotFound
	}

	usr, err = gh.Repositories(req.r.GetRuntime()).User.GetByID(
		userID,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)

	return usr, err
}
