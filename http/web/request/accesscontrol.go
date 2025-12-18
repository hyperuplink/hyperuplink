package request

import (
	"slices"

	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/services/repositories/common"
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

	usr, err = req.r.GetRuntime().Repositories.User.GetByID(
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
