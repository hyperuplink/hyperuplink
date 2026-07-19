package settings

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	Timezone string `json:"timezone" form:"timezone" validate:"required,timezone"`
}

func Update(
	rt *runtime.Runtime,
	usr *user.User,
	in *UpdateInput,
) (err error) {
	usr.Timezone = in.Timezone

	return rt.Repositories.User.Update(usr)
}

func viewToggle(profile *setting.UserProfile, view string) *bool {
	switch view {
	case setting.UserProfileViewBanner:
		return &profile.ShowBanner
	case setting.UserProfileViewFooter:
		return &profile.ShowFooter
	case setting.UserProfileViewProfilePictures:
		return &profile.ShowProfilePictures
	}

	return nil
}

func ToggleView(
	rt *runtime.Runtime,
	usr *user.User,
	view string,
) (enabled bool, err error) {
	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(
		rt.Repositories.Setting,
		usr.ID,
	)
	if err != nil {
		return false, err
	}

	toggle := viewToggle(&settingUserProfile.JSONValue, view)
	if toggle == nil {
		return false, errs.ErrValidation
	}
	*toggle = !*toggle

	if err = settingRepo.Update[setting.UserProfile](
		rt.Repositories.Setting,
		settingUserProfile,
	); err != nil {
		return false, err
	}

	return *toggle, nil
}
