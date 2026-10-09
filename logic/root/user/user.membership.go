package user

import (
	"slices"

	"github.com/google/uuid"

	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type MembershipInput struct {
	MemberOf []string `json:"member_of" form:"member_of"`
}

func UpdateMembership(
	rt *runtime.Runtime,
	username string,
	in *MembershipInput,
) (err error) {
	var usr *user.User
	if id, perr := uuid.Parse(username); perr == nil {
		usr, err = gh.Repositories(rt).User.GetByUUID(id, common.QueryOptions{Limit: 1})
	} else {
		usr, err = gh.Repositories(rt).User.GetByUsername(username, common.QueryOptions{Limit: 1})
	}
	if err != nil {
		return err
	}

	groups, err := gh.Repositories(rt).Group.All(common.QueryOptions{})
	if err != nil {
		return err
	}

	memberOf := []string{}
	for _, grp := range *groups {
		if slices.Contains(in.MemberOf, grp.ID) {
			memberOf = append(memberOf, grp.ID)
		}
	}
	usr.MemberOf = memberOf

	return gh.Repositories(rt).User.Update(usr)
}
