package topics

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type VoteInput struct {
	Topic     *vtopic.VTopic
	AuthorID  uuid.UUID
	Selection int
}

type Voted struct {
	Topic *vtopic.VTopic
	Poll  *logictopics.Poll
}

func Vote(
	rt *runtime.Runtime,
	in *VoteInput,
	perms *permission.Resolution,
) (voted *Voted, err error) {
	top := in.Topic

	if !perms.CanWriteSlug(top.CategorySlug) {
		return nil, errs.ErrForbidden
	}

	if err = logictopics.PollVote(rt, &logictopics.PollVoteInput{
		Topic:     &top.Topic,
		AuthorID:  in.AuthorID,
		Selection: in.Selection,
	}); err != nil {
		return nil, err
	}

	poll, err := logictopics.PollView(rt, &logictopics.PollViewInput{
		Topic:    &top.Topic,
		ViewerID: uuid.NullUUID{UUID: in.AuthorID, Valid: true},
		CanWrite: true,
	})
	if err != nil {
		return nil, err
	}

	return &Voted{Topic: top, Poll: poll}, nil
}
