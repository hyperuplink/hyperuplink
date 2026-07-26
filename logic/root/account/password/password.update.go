package password

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type UpdateInput struct {
	CurrentPassword   string `json:"current_password" form:"current_password" validate:"required,min=8,max=64"`
	NewPassword       string `json:"new_password" form:"new_password" validate:"required,min=8,max=64"`
	NewPasswordRepeat string `json:"new_password_repeat" form:"new_password_repeat" validate:"required,eqcsfield=NewPassword"`
	OTPCode           string `json:"otp_code" form:"otp_code" validate:"omitempty,numeric,len=6"`
}

func Update(
	rt *runtime.Runtime,
	usr *user.User,
	in *UpdateInput,
) (err error) {
	var match bool
	if match, _, err = usr.CheckPassword(in.CurrentPassword); !match {
		if err == nil {
			err = errs.ErrPasswordWrong
		}
	}
	if err != nil {
		return err
	}

	if usr.OTPEnabled && !user.ValidateOTP(usr.OTPSecret, in.OTPCode) {
		return errs.ErrOTPCodeWrong
	}

	if err = usr.SetPassword(in.NewPassword); err != nil {
		return err
	}

	return gh.Repositories(rt).User.Update(usr)
}
