package settings

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	Language string `json:"language" form:"language" validate:"required,oneof=de en es fr it ro"`
	Timezone string `json:"timezone" form:"timezone" validate:"required,timezone"`
}

func Update(
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
	in *UpdateInput,
) (err error) {
	usr.Language = in.Language
	usr.Timezone = in.Timezone

	return gh.Repositories(rt).User.Update(ctx, usr)
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
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
	view string,
) (enabled bool, err error) {
	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(ctx,
		gh.Repositories(rt).Setting,
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

	if err = settingRepo.Update[setting.UserProfile](ctx,
		gh.Repositories(rt).Setting,
		settingUserProfile,
	); err != nil {
		return false, err
	}

	return *toggle, nil
}
