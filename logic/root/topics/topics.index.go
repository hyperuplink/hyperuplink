package topics

import (
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

type Input struct {
	ForumID  uuid.NullUUID
	Page     int
	PerPage  int
	ViewerID uuid.NullUUID
}

type View struct {
	Forum  *vforum.VForum   `json:"forum,omitempty"`
	Topics *[]vtopic.VTopic `json:"topics"`
	Total  int64            `json:"total"`
	Pages  int              `json:"pages"`
	Unread map[string]bool  `json:"unread,omitempty"`
}

func Index(
	rt *runtime.Runtime,
	in *Input,
	perms *permission.Resolution,
) (view *View, err error) {
	qo := common.QueryOptions{
		OrderBy: "updated_at",
		Order:   common.Descending,
		Limit:   in.PerPage,
		Page:    in.Page,
	}

	view = new(View)

	var tops *[]vtopic.VTopic
	var total int64

	if in.ForumID.Valid {
		var fum *vforum.VForum
		fum, err = gh.Repositories(rt).Forum.VGetByUUID(
			in.ForumID.UUID,
			common.QueryOptions{Limit: 1},
		)
		if err != nil {
			return nil, err
		}

		if !perms.CanReadID(fum.CategoryID) {
			return nil, errs.ErrForbidden
		}

		tops, total, err = gh.Repositories(rt).Topic.VAllForForumUUID(fum.ID, qo)
		if err != nil {
			return nil, err
		}

		view.Forum = fum
	} else {
		tops, total, err = gh.Repositories(rt).Topic.VAllForReadableSlugs(
			perms.AllowedReadSlugs(),
			qo,
		)
		if err != nil {
			return nil, err
		}
	}

	view.Topics = tops
	view.Total = total
	view.Pages = paging.Pages(total, in.PerPage)

	if in.ViewerID.Valid {
		view.Unread, err = logicactivity.UnreadTopics(rt, in.ViewerID.UUID, tops)
		if err != nil {
			return nil, err
		}
	}

	return view, nil
}
