package activity

import (
	"time"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

func RecordTopicView(
	rt *runtime.Runtime,
	actorID uuid.UUID,
	topicID uuid.UUID,
) {
	if err := gh.Activity(rt).Record(
		activity.NewTopicView(actorID, topicID),
	); err != nil {
		rt.Error("error", err)
	}
}

func UnreadTopics(
	rt *runtime.Runtime,
	actorID uuid.UUID,
	topics *[]vtopic.VTopic,
) (unread map[string]bool, err error) {
	unread = make(map[string]bool)

	if topics == nil || len(*topics) == 0 {
		return unread, nil
	}

	subjectIDs := make([]uuid.UUID, 0, len(*topics))
	for _, top := range *topics {
		subjectIDs = append(subjectIDs, top.ID)
	}

	var markers map[uuid.UUID]time.Time
	markers, err = gh.Repositories(rt).Activity.MarkersForActor(
		activity.TopicView,
		actorID,
		subjectIDs,
		common.QueryOptions{},
	)
	if err != nil {
		return nil, err
	}

	for _, top := range *topics {
		readAt, ok := markers[top.ID]
		if !ok {
			unread[top.ID.String()] = true
			continue
		}

		latest := top.CreatedAt
		if top.LastReplyAt.Valid {
			latest = top.LastReplyAt
		}

		unread[top.ID.String()] = latest.Valid && latest.Time.After(readAt)
	}

	return unread, nil
}
