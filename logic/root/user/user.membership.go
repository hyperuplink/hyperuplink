package user

import (
	"slices"

	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type MembershipInput struct {
	MemberOf []string `json:"member_of" form:"member_of"`
}

func UpdateMembership(
	rt *runtime.Runtime,
	username string,
	in *MembershipInput,
) (err error) {
	usr, err := rt.Repositories.User.GetByUsername(
		username,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		return err
	}

	groups, err := rt.Repositories.Group.All(common.QueryOptions{})
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

	return rt.Repositories.User.Update(usr)
}
