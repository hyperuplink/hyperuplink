package topics

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type ShowInput struct {
	ForumSlug string
	TopicSlug string
	Page      int
	PerPage   int
	ViewerID  uuid.NullUUID
}

type View struct {
	Topic   *vtopic.VTopic   `json:"topic"`
	Poll    *Poll            `json:"poll,omitempty"`
	Replies *[]vreply.VReply `json:"replies"`
	Total   int64            `json:"total"`
	Pages   int              `json:"pages"`
}

func Show(
	rt *runtime.Runtime,
	in *ShowInput,
	perms *permission.Resolution,
) (view *View, err error) {
	top, err := rt.Repositories.Topic.VGetBySlugs(
		in.ForumSlug,
		in.TopicSlug,
		common.QueryOptions{
			Limit: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	if !perms.CanReadSlug(top.CategorySlug) {
		return nil, errs.ErrForbidden
	}

	if in.ViewerID.Valid {
		logicactivity.RecordTopicView(rt, in.ViewerID.UUID, top.ID)
	}

	poll, err := PollView(rt, &PollViewInput{
		Topic:    &top.Topic,
		ViewerID: in.ViewerID,
		CanWrite: perms.CanWriteSlug(top.CategorySlug),
	})
	if err != nil {
		return nil, err
	}

	limit := in.PerPage
	offAdjust := 1

	// If we are on the first page, we subtract 1 from limit due to the Topic,
	// and we set offadjust to 0 because we only need to adjust the offset for
	// all pages past the first one.
	if in.Page == 1 {
		limit -= 1
		offAdjust = 0
	}
	reps, total, err := rt.Repositories.Reply.VAllForTopicUUID(
		top.ID,
		common.QueryOptions{
			OrderBy:   "created_at",
			Order:     common.Ascending,
			Limit:     limit,
			Page:      in.Page,
			OffAdjust: offAdjust,
		},
	)
	if err != nil {
		return nil, err
	}
	// We add the Topic to the total
	total += 1

	return &View{
		Topic:   top,
		Poll:    poll,
		Replies: reps,
		Total:   total,
		Pages:   paging.Pages(total, in.PerPage),
	}, nil
}
