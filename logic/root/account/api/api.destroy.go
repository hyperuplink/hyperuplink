package api

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type DestroyInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}

func Destroy(
	rt *runtime.Runtime,
	userID uuid.UUID,
	in *DestroyInput,
) (err error) {
	var id uuid.UUID
	if id, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	key, err := rt.Repositories.APIKey.GetByUUID(id, common.QueryOptions{Limit: 1})
	if err != nil {
		return err
	}

	if key.UserID != userID {
		return errs.ErrNoRows
	}

	return rt.Repositories.APIKey.Delete(key)
}
