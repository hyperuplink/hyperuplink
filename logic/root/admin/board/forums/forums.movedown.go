package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func MoveDown(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	fum := new(forum.Forum)
	if fum.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	if _, err = gh.Repositories(rt).Forum.GetByUUID(fum.ID, common.QueryOptions{Limit: 1}); err != nil {
		return err
	}

	return gh.Repositories(rt).Forum.MoveDown(fum)
}
