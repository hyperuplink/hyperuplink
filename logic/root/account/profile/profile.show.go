package profile

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(
	rt *runtime.Runtime,
	usr *user.User,
) (view *View, err error) {
	settingProfiles, err := settingRepo.GetByID[setting.Profiles](
		rt.Repositories.Setting,
		"profiles",
	)
	if err != nil {
		return nil, err
	}

	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(
		rt.Repositories.Setting,
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
