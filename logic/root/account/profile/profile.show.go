package profile

import (
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(
	rt *runtime.Runtime,
	usr *user.User,
) (view *View, err error) {
	settingProfiles, err := settingRepo.GetByID[setting.Profiles](
		gh.Repositories(rt).Setting,
		"profiles",
	)
	if err != nil {
		return nil, err
	}

	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(
		gh.Repositories(rt).Setting,
		usr.ID,
	)
	if err != nil {
		return nil, err
	}

	return &View{
		Profiles:    &settingProfiles.JSONValue,
		UserProfile: &settingUserProfile.JSONValue,
	}, nil
}
