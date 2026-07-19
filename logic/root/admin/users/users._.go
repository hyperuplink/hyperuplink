package users

import (
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type UserInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}

func getUser(rt *runtime.Runtime, id string) (*user.User, error) {
	return rt.Repositories.User.GetByID(
		id,
		common.QueryOptions{
			WithBanned:  true,
			WithDeleted: true,
			Limit:       1,
		},
	)
}
