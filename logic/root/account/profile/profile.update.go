package profile

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Update(
	rt *runtime.Runtime,
	usr *user.User,
	in *UpdateInput,
) (err error) {
	if in.ProfilePicture != nil && in.ProfilePicture.Filename != "" {
		if err = StorePicture(rt, usr, in.ProfilePicture); err != nil {
			return err
		}
	}

	usr.SignatureText = in.SignatureText
	if usr.SignatureText == "" {
		usr.SignatureHTML = ""
	} else {
		if usr.SignatureHTML, err = rt.Markdown.Convert(usr.SignatureText); err != nil {
			return err
		}
	}

	if err = rt.Repositories.User.Update(usr); err != nil {
		return err
	}

	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(
		rt.Repositories.Setting,
		usr.ID,
	)
	if err != nil {
		return err
	}

	settingUserProfile.JSONValue.NotifyOnReply = in.NotifyOnReply

	return settingRepo.Update[setting.UserProfile](
		rt.Repositories.Setting,
		settingUserProfile,
	)
}
