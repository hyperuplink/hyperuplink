package forums

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type ShowInput struct {
	ForumSlug string
	Page      int
	PerPage   int
	ViewerID  uuid.NullUUID
}

type View struct {
	Forum  *vforum.VForum  `json:"forum"`
	Topics []vtopic.VTopic `json:"topics"`
	Total  int64           `json:"total"`
	Pages  int             `json:"pages"`
	Unread map[string]bool `json:"unread,omitempty"`
}

func Show(
	ctx context.Context,
	rt *runtime.Runtime,
	in *ShowInput,
	perms *permission.Resolution,
) (view *View, err error) {
	fum, err := gh.Repositories(rt).Forum.VGetBySlug(ctx,
		in.ForumSlug,
		common.QueryOptions{
			Limit: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	if !perms.CanReadID(fum.CategoryID) {
		return nil, errs.ErrForbidden
	}

	tops, total, err := gh.Repositories(rt).Topic.VAllForForumUUID(ctx,
		fum.ID,
		common.QueryOptions{
			OrderBy: "last_activity_at",
			Order:   common.Descending,
			Limit:   in.PerPage,
			Page:    in.Page,
		},
	)
	if err != nil {
		return nil, err
	}

	view = &View{
		Forum:  fum,
		Topics: tops,
		Total:  total,
		Pages:  paging.Pages(total, in.PerPage),
	}

	if in.ViewerID.Valid {
		view.Unread, err = logicactivity.UnreadTopics(ctx, rt, in.ViewerID.UUID, tops)
		if err != nil {
			return nil, err
		}
	}

	return view, nil
}
