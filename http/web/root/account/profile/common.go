package profile

import (
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (r *Route) getUser(req *request.Request) (usr *user.User, err error) {
	var userID string
	var ok bool
	if userID, ok = req.Session.GetUserID(); !ok {
		return nil, errs.ErrUserIDNotFound
	}

	usr, err = r.Runtime.Repositories.User.GetByID(
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
