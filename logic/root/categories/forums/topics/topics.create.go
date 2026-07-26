package topics

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

type CreateReplyInput struct {
	Text    string `json:"text" form:"text" validate:"required,min=1"`
	TopicID string `json:"topic_id" form:"topic_id" validate:"required,uuid"`
	ReplyID string `json:"reply_id" form:"reply_id" validate:"omitempty,uuid"`
}

type CreatedReply struct {
	Reply *reply.Reply
	Topic *topic.Topic
	Forum *forum.Forum
}

func CreateReply(
	rt *runtime.Runtime,
	authorID uuid.UUID,
	perms *permission.Resolution,
	in *CreateReplyInput,
	attachments func(authorID uuid.UUID) ([]uuid.UUID, error),
) (created *CreatedReply, err error) {
	rep := new(reply.Reply)
	rep.ShortID = shortuuid.New()
	if rep.TopicID, err = uuid.Parse(in.TopicID); err != nil {
		return nil, err
	}

	top, err := gh.Repositories(rt).Topic.GetByUUID(
		rep.TopicID,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return nil, err
	}
	fum, err := gh.Repositories(rt).Forum.GetByUUID(
		top.ForumID,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return nil, err
	}
	if !perms.CanWriteID(fum.CategoryID) {
		return nil, errs.ErrForbidden
	}

	if in.ReplyID != "" {
		var replyUUID uuid.UUID
		if replyUUID, err = uuid.Parse(in.ReplyID); err != nil {
			return nil, err
		}
		rep.ReplyID = uuid.NullUUID{
			UUID:  replyUUID,
			Valid: true,
		}
	} else {
		rep.ReplyID = uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}
	rep.AuthorID = authorID
	rep.Text = in.Text
	if rep.HTML, err = rt.Markdown().Convert(rep.Text); err != nil {
		return nil, err
	}

	if attachments != nil {
		if rep.AttachmentIDs, err = attachments(authorID); err != nil {
			return nil, err
		}
	}

	if rep.ID, err = gh.Repositories(rt).Reply.Create(rep); err != nil {
		return nil, err
	}

	return &CreatedReply{
		Reply: rep,
		Topic: top,
		Forum: fum,
	}, nil
}

func ReplyPages(
	rt *runtime.Runtime,
	topicID uuid.UUID,
	perPage int,
) (pages int, err error) {
	total, err := gh.Repositories(rt).Reply.VAllCountForTopicUUID(
		topicID,
		common.QueryOptions{},
	)
	if err != nil {
		return 0, err
	}

	return paging.Pages(total, perPage), nil
}
