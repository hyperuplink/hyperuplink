package session

import (
	"context"
	"errors"
	"sync"

	"xn--gckvb8fzb.com/glides/helpers/password"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type SignInInput struct {
	Username string `json:"username" form:"username" validate:"required,min=2,max=32"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=64"`
	OTPCode  string `json:"otp_code" form:"otp_code" validate:"omitempty,numeric,len=6"`
}

var (
	decoyOnce sync.Once
	decoyHash string
)

func decoy() string {
	decoyOnce.Do(func() {
		decoyHash, _ = password.HashRandom()
	})

	return decoyHash
}

func CheckPassword(
	ctx context.Context,
	rt *runtime.Runtime,
	username string,
	pw string,
) (usr *user.User, err error) {
	usr, err = gh.Repositories(rt).User.GetByUsername(ctx,
		username,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			password.Check(decoy(), pw)
			return nil, errs.ErrUsernamePasswordWrong
		}
		return nil, err
	}

	match, _, cerr := usr.CheckPassword(pw)
	if cerr != nil || !match {
		return nil, errs.ErrUsernamePasswordWrong
	}

	return usr, nil
}

func SignIn(
	ctx context.Context,
	rt *runtime.Runtime,
	in *SignInInput,
) (usr *user.User, err error) {
	if usr, err = CheckPassword(ctx, rt, in.Username, in.Password); err != nil {
		return nil, err
	}

	if usr.OTPEnabled {
		if in.OTPCode == "" {
			return nil, errs.ErrOTPRequired
		}
		if !user.ValidateOTP(usr.OTPSecret, in.OTPCode) {
			return nil, errs.ErrOTPCodeWrong
		}
	}

	return usr, nil
}
