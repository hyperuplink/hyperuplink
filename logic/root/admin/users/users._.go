package users

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type UserInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}

func getUser(rt *runtime.Runtime, id string) (*user.User, error) {
	return gh.Repositories(rt).User.GetByID(
		id,
		common.QueryOptions{
			WithBanned:  true,
			WithDeleted: true,
			Limit:       1,
		},
	)
}
