package report

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

type ResolvedPost struct {
	Target   postevent.PostEventTarget
	Text     string
	TopicID  uuid.NullUUID
	ReplyID  uuid.NullUUID
	AuthorID uuid.UUID
}

type CreateInput struct {
	Target     string `json:"target" form:"target" validate:"required,oneof=topic reply"`
	ID         string `json:"id" form:"id" validate:"required"`
	ReportType int    `json:"report_type" form:"report_type" validate:"oneof=0 1 2"`
}

func ResolvePost(
	rt *runtime.Runtime,
	target string,
	shortID string,
) (post *ResolvedPost, err error) {
	id, perr := uuid.Parse(shortID)

	switch postevent.PostEventTarget(target) {
	case postevent.Topic:
		var top *topic.Topic
		if perr == nil {
			top, err = gh.Repositories(rt).Topic.GetByUUID(id, common.QueryOptions{Limit: 1})
		} else {
			top, err = gh.Repositories(rt).Topic.GetByShortID(shortID, common.QueryOptions{Limit: 1})
		}
		if err != nil {
			return nil, err
		}
		return &ResolvedPost{
			Target:   postevent.Topic,
			Text:     top.Text,
			TopicID:  uuid.NullUUID{UUID: top.ID, Valid: true},
			AuthorID: top.AuthorID,
		}, nil
	case postevent.Reply:
		var rep *reply.Reply
		if perr == nil {
			rep, err = gh.Repositories(rt).Reply.GetByUUID(id, common.QueryOptions{Limit: 1})
		} else {
			rep, err = gh.Repositories(rt).Reply.GetByShortID(shortID, common.QueryOptions{Limit: 1})
		}
		if err != nil {
			return nil, err
		}
		return &ResolvedPost{
			Target:   postevent.Reply,
			Text:     rep.Text,
			ReplyID:  uuid.NullUUID{UUID: rep.ID, Valid: true},
			AuthorID: rep.AuthorID,
		}, nil
	default:
		return nil, errs.ErrTargetIDNotFound
	}
}
