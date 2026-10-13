package topics

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
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type ShowInput struct {
	CategorySlug string
	ForumSlug    string
	TopicSlug    string
	Page         int
	PerPage      int
	ViewerID     uuid.NullUUID
}

type View struct {
	Topic   *vtopic.VTopic  `json:"topic"`
	Poll    *Poll           `json:"poll,omitempty"`
	Replies []vreply.VReply `json:"replies"`
	Total   int64           `json:"total"`
	Pages   int             `json:"pages"`
}

func Show(
	ctx context.Context,
	rt *runtime.Runtime,
	in *ShowInput,
	perms *permission.Resolution,
) (view *View, err error) {
	top, err := gh.Repositories(rt).Topic.VGetBySlugs(ctx,
		in.CategorySlug,
		in.ForumSlug,
		in.TopicSlug,
		common.QueryOptions{
			Limit: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	return viewForTopic(ctx, rt, top, in.Page, in.PerPage, in.ViewerID, perms)
}

func viewForTopic(
	ctx context.Context,
	rt *runtime.Runtime,
	top *vtopic.VTopic,
	page int,
	perPage int,
	viewerID uuid.NullUUID,
	perms *permission.Resolution,
) (view *View, err error) {
	if !perms.CanReadSlug(top.CategorySlug) {
		return nil, errs.ErrForbidden
	}

	if viewerID.Valid {
		logicactivity.RecordTopicView(ctx, rt, viewerID.UUID, top.ID)
	}

	poll, err := PollView(ctx, rt, &PollViewInput{
		Topic:    &top.Topic,
		ViewerID: viewerID,
		CanWrite: perms.CanWriteSlug(top.CategorySlug),
	})
	if err != nil {
		return nil, err
	}

	limit := perPage
	offAdjust := 1
	if page == 1 {
		limit -= 1
		offAdjust = 0
	}
	reps, total, err := gh.Repositories(rt).Reply.VAllForTopicUUID(ctx,
		top.ID,
		common.QueryOptions{
			Limit:     limit,
			Page:      page,
			OffAdjust: offAdjust,
		},
	)
	if err != nil {
		return nil, err
	}
	total += 1

	return &View{
		Topic:   top,
		Poll:    poll,
		Replies: reps,
		Total:   total,
		Pages:   paging.Pages(total, perPage),
	}, nil
}
