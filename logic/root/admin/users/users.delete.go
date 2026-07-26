package users

import (
	"time"

	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

func Delete(rt *runtime.Runtime, in *UserInput) (err error) {
	usr, err := getUser(rt, in.ID)
	if err != nil {
		return err
	}

	usr.SetDeletedAt(time.Now())

	return gh.Repositories(rt).User.Update(usr)
}
