package account

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type View struct {
	Groups  []group.Group
	Profile *setting.UserProfile
}

func Show(
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
) (view *View, err error) {
	view = &View{Groups: []group.Group{}}

	if len(usr.MemberOf) > 0 {
		var groups []group.Group
		if groups, err = gh.Repositories(rt).Group.AllByIDs(ctx, usr.MemberOf, common.QueryOptions{
			OrderBy: "name",
			Order:   common.Ascending,
		}); err != nil {
			return nil, err
		}
		view.Groups = groups
	}

	profile, err := settingRepo.GetOrCreateUserProfile(ctx, gh.Repositories(rt).Setting, usr.ID)
	if err != nil {
		return nil, err
	}
	view.Profile = &profile.JSONValue

	return view, nil
}
