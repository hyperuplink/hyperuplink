package report

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type resolvedPost struct {
	Target   postevent.PostEventTarget
	Text     string
	TopicID  uuid.NullUUID
	ReplyID  uuid.NullUUID
	AuthorID uuid.UUID
}

func (r *Route) resolvePost(target, shortID string) (*resolvedPost, error) {
	switch postevent.PostEventTarget(target) {
	case postevent.Topic:
		top, err := r.Runtime.Repositories.Topic.GetByShortID(
			shortID,
			common.QueryOptions{Limit: 1},
		)
		if err != nil {
			return nil, err
		}
		return &resolvedPost{
			Target:   postevent.Topic,
			Text:     top.Text,
			TopicID:  uuid.NullUUID{UUID: top.ID, Valid: true},
			AuthorID: top.AuthorID,
		}, nil
	case postevent.Reply:
		rep, err := r.Runtime.Repositories.Reply.GetByShortID(
			shortID,
			common.QueryOptions{Limit: 1},
		)
		if err != nil {
			return nil, err
		}
		return &resolvedPost{
			Target:   postevent.Reply,
			Text:     rep.Text,
			ReplyID:  uuid.NullUUID{UUID: rep.ID, Valid: true},
			AuthorID: rep.AuthorID,
		}, nil
	default:
		return nil, errs.ErrTargetIDNotFound
	}
}
