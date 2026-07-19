package users

import (
	"time"

	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Delete(rt *runtime.Runtime, in *UserInput) (err error) {
	usr, err := getUser(rt, in.ID)
	if err != nil {
		return err
	}

	usr.SetDeletedAt(time.Now())

	return rt.Repositories.User.Update(usr)
}
