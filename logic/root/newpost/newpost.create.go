package newpost

import (
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type CreateInput struct {
	Name        string   `json:"name" form:"name" validate:"required,min=1,max=78"`
	Text        string   `json:"text" form:"text" validate:"required,min=1"`
	ForumID     string   `json:"forum_id" form:"forum_id" validate:"required,uuid"`
	Kind        string   `json:"kind" form:"kind" validate:"omitempty,oneof=regular poll"`
	PollOptions []string `json:"poll_options" form:"poll_options" validate:"omitempty,dive,max=78"`
	PollEndsAt  string   `json:"poll_ends_at" form:"poll_ends_at"`
}

func Create(
	rt *runtime.Runtime,
	authorID uuid.UUID,
	perms *permission.Resolution,
	location *time.Location,
	in *CreateInput,
	attachments func(authorID uuid.UUID) ([]uuid.UUID, error),
) (vtop *vtopic.VTopic, err error) {
	top := new(topic.Topic)
	top.ShortID = shortuuid.New()
	top.Name = in.Name
	top.SetSlugFromName()
	if top.ForumID, err = uuid.Parse(in.ForumID); err != nil {
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

	top.AuthorID = authorID
	top.Kind = topic.Regular
	top.Anonymous = false
	top.Text = in.Text
	if top.HTML, err = rt.Markdown().Convert(top.Text); err != nil {
		return nil, err
	}

	if topic.Kind(in.Kind) == topic.Poll {
		if err = ApplyPoll(rt, top, &PollInput{
			Options:  in.PollOptions,
			EndsAt:   in.PollEndsAt,
			Location: location,
		}); err != nil {
			return nil, err
		}
	}

	if attachments != nil {
		if top.AttachmentIDs, err = attachments(authorID); err != nil {
			return nil, err
		}
	}

	if _, err = gh.Repositories(rt).Topic.Create(top); err != nil {
		return nil, err
	}

	return gh.Repositories(rt).Topic.VGetByForumUUIDSlug(
		top.ForumID,
		top.Slug,
		common.QueryOptions{
			Limit: 1,
		},
	)
}
