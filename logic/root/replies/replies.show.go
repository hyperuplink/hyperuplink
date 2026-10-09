package replies

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type Shown struct {
	Reply    *vreply.VReply
	Topic    *vtopic.VTopic
	Position int64
}

func Resolve(
	rt *runtime.Runtime,
	idOrShortID string,
) (rep *vreply.VReply, err error) {
	if id, perr := uuid.Parse(idOrShortID); perr == nil {
		return gh.Repositories(rt).Reply.VGetByUUID(id, common.QueryOptions{Limit: 1})
	}

	return gh.Repositories(rt).Reply.VGetByShortID(idOrShortID, common.QueryOptions{Limit: 1})
}

func Show(
	rt *runtime.Runtime,
	idOrShortID string,
	perms *permission.Resolution,
) (shown *Shown, err error) {
	rep, err := Resolve(rt, idOrShortID)
	if err != nil {
		return nil, err
	}

	top, err := gh.Repositories(rt).Topic.VGetByUUID(rep.TopicID, common.QueryOptions{Limit: 1})
	if err != nil {
		return nil, err
	}

	if !perms.CanReadSlug(top.CategorySlug) {
		return nil, errs.ErrForbidden
	}

	position, err := gh.Repositories(rt).Reply.PositionForUUID(rep.ID)
	if err != nil {
		return nil, err
	}

	return &Shown{Reply: rep, Topic: top, Position: position}, nil
}
