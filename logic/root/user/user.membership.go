package user

import (
	"slices"

	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

type MembershipInput struct {
	MemberOf []string `json:"member_of" form:"member_of"`
}

func UpdateMembership(
	rt *runtime.Runtime,
	username string,
	in *MembershipInput,
) (err error) {
	usr, err := gh.Repositories(rt).User.GetByUsername(
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
