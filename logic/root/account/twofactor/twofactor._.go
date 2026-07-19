package twofactor

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Enable(
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

	return rt.Repositories.User.Update(usr)
}

func Disable(
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

	return rt.Repositories.User.Update(usr)
}
