package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func MoveUp(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	fum := new(forum.Forum)
	if fum.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return rt.Repositories.Forum.MoveUp(fum)
}
