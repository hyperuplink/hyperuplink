package twofactor

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func Enable(
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
	secret string,
	code string,
) (err error) {
	if usr.OTPEnabled {
		return nil
	}

	if !user.ValidateOTP(secret, code) {
		return errs.ErrOTPCodeWrong
	}

	usr.EnableOTP(secret)

	return gh.Repositories(rt).User.Update(ctx, usr)
}

func Disable(
	ctx context.Context,
	rt *runtime.Runtime,
	usr *user.User,
	currentPassword string,
) (err error) {
	if !usr.OTPEnabled {
		return nil
	}

	var match bool
	if match, _, err = usr.CheckPassword(currentPassword); !match {
		if err == nil {
			err = errs.ErrPasswordWrong
		}
	}
	if err != nil {
		return err
	}

	usr.DisableOTP()

	return gh.Repositories(rt).User.Update(ctx, usr)
}
