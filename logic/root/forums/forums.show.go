package forums

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type TopicsInput struct {
	Forum    *vforum.VForum
	Sort     topic.Sort
	Page     int
	PerPage  int
	ViewerID uuid.NullUUID
}

type Topics struct {
	Topics []vtopic.VTopic
	Total  int64
	Unread map[string]bool
}

func Resolve(
	ctx context.Context,
	rt *runtime.Runtime,
	idOrSlug string,
) (fum *vforum.VForum, err error) {
	if id, perr := uuid.Parse(idOrSlug); perr == nil {
		return gh.Repositories(rt).Forum.VGetByUUID(ctx, id, common.QueryOptions{Limit: 1})
	}

	return gh.Repositories(rt).Forum.VGetBySlug(ctx, idOrSlug, common.QueryOptions{Limit: 1})
}

func Show(
	ctx context.Context,
	rt *runtime.Runtime,
	idOrSlug string,
	perms *permission.Resolution,
) (fum *vforum.VForum, err error) {
	if fum, err = Resolve(ctx, rt, idOrSlug); err != nil {
		return nil, err
	}

	if !perms.CanReadID(fum.CategoryID) {
		return nil, errs.ErrForbidden
	}

	return fum, nil
}

func ListTopics(
	ctx context.Context,
	rt *runtime.Runtime,
	in *TopicsInput,
) (page *Topics, err error) {
	page = new(Topics)

	if page.Topics, page.Total, err = gh.Repositories(rt).Topic.VAllForForumUUIDSorted(ctx,
		in.Forum.ID,
		in.Sort,
		in.PerPage,
		in.Page,
	); err != nil {
		return nil, err
	}

	if in.ViewerID.Valid {
		if page.Unread, err = logicactivity.UnreadTopics(ctx, rt, in.ViewerID.UUID, page.Topics); err != nil {
			return nil, err
		}
	}

	return page, nil
}
