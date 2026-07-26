package api

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func List(
	rt *runtime.Runtime,
	userID uuid.UUID,
) (keys *[]apikey.APIKey, err error) {
	return gh.Repositories(rt).APIKey.AllForUserUUID(
		userID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
		},
	)
}
