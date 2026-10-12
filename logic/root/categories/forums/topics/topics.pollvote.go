package topics

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

type PollVoteBySlugsInput struct {
	ForumSlug string
	TopicSlug string
	AuthorID  uuid.UUID
	Selection int
}

func PollVoteBySlugs(
	ctx context.Context,
	rt *runtime.Runtime,
	perms *permission.Resolution,
	in *PollVoteBySlugsInput,
) (err error) {
	vtop, err := gh.Repositories(rt).Topic.VGetBySlugs(ctx,
		in.ForumSlug,
		in.TopicSlug,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return err
	}

	if !perms.CanWriteSlug(vtop.CategorySlug) {
		return errs.ErrForbidden
	}

	top, err := gh.Repositories(rt).Topic.GetByUUID(ctx,
		vtop.ID,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return err
	}

	return PollVote(ctx, rt, &PollVoteInput{
		Topic:     top,
		AuthorID:  in.AuthorID,
		Selection: in.Selection,
	})
}
