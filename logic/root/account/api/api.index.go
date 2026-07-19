package api

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func List(
	rt *runtime.Runtime,
	userID uuid.UUID,
) (keys *[]apikey.APIKey, err error) {
	return rt.Repositories.APIKey.AllForUserUUID(
		userID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
		},
	)
}
