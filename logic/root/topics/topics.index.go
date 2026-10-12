package topics

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type Input struct {
	Sort     topic.Sort
	Page     int
	PerPage  int
	ViewerID uuid.NullUUID
}

type Feed struct {
	Topics []vtopic.VTopic
	Total  int64
	Unread map[string]bool
}

func Index(
	rt *runtime.Runtime,
	in *Input,
	perms *permission.Resolution,
) (feed *Feed, err error) {
	feed = new(Feed)

	if feed.Topics, feed.Total, err = gh.Repositories(rt).Topic.VAllForReadableSlugsSorted(
		perms.AllowedReadSlugs(),
		in.Sort,
		in.PerPage,
		in.Page,
	); err != nil {
		return nil, err
	}

	if in.ViewerID.Valid {
		if feed.Unread, err = logicactivity.UnreadTopics(rt, in.ViewerID.UUID, feed.Topics); err != nil {
			return nil, err
		}
	}

	return feed, nil
}
