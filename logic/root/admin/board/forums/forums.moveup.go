package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func MoveUp(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	fum := new(forum.Forum)
	if fum.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return gh.Repositories(rt).Forum.MoveUp(fum)
}
