package users

import (
	"context"
	"time"

	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

func Unban(ctx context.Context, rt *runtime.Runtime, in *UserInput) (err error) {
	usr, err := getUser(ctx, rt, in.ID)
	if err != nil {
		return err
	}

	usr.SetBannedAt(time.Time{})

	return gh.Repositories(rt).User.Update(ctx, usr)
}
