package topics

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type PollVoteBySlugsInput struct {
	ForumSlug string
	TopicSlug string
	AuthorID  uuid.UUID
	Selection int
}

func PollVoteBySlugs(
	rt *runtime.Runtime,
	perms *permission.Resolution,
	in *PollVoteBySlugsInput,
) (err error) {
	vtop, err := rt.Repositories.Topic.VGetBySlugs(
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

	top, err := rt.Repositories.Topic.GetByUUID(
		vtop.ID,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return err
	}

	return PollVote(rt, &PollVoteInput{
		Topic:     top,
		AuthorID:  in.AuthorID,
		Selection: in.Selection,
	})
}
