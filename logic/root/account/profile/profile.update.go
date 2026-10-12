package profile

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/markdown"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Update(
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
	in *UpdateInput,
) (err error) {
	if in.ProfilePicture != nil && in.ProfilePicture.Filename != "" {
		if err = StorePicture(ctx, rt, usr, in.ProfilePicture); err != nil {
			return err
		}
	}

	usr.SignatureText = in.SignatureText
	if usr.SignatureText == "" {
		usr.SignatureHTML = ""
	} else {
		if usr.SignatureHTML, err = markdown.From(rt).Convert(usr.SignatureText); err != nil {
			return err
		}
	}

	if err = gh.Repositories(rt).User.Update(ctx, usr); err != nil {
		return err
	}

	settingUserProfile, err := settingRepo.GetOrCreateUserProfile(ctx,
		gh.Repositories(rt).Setting,
		usr.ID,
	)
	if err != nil {
		return err
	}

	settingUserProfile.JSONValue.NotifyOnReply = in.NotifyOnReply

	return settingRepo.Update[setting.UserProfile](ctx,
		gh.Repositories(rt).Setting,
		settingUserProfile,
	)
}
