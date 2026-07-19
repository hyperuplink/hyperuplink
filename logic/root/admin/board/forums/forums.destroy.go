package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Destroy(
	rt *runtime.Runtime,
	in *DestroyInput,
) (err error) {
	fum := new(forum.Forum)
	if fum.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return rt.Repositories.Forum.Delete(fum)
}
