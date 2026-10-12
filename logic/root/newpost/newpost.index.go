package newpost

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type CategoryWithForums struct {
	Category category.Category `json:"category"`
	Forums   []forum.Forum     `json:"forums"`
}

type FormViewInput struct {
	ForumSlug string
	TopicSlug string
	ReplyID   string
}

type FormView struct {
	CategoriesForums []CategoryWithForums `json:"categories_forums"`
	AllowPoll        bool                 `json:"allow_poll"`
	PollOptionsMax   int                  `json:"poll_options_max"`
	Forum            *vforum.VForum       `json:"forum,omitempty"`
	Topic            *vtopic.VTopic       `json:"topic,omitempty"`
	Reply            *reply.Reply         `json:"reply,omitempty"`
}

func View(
	ctx context.Context,
	rt *runtime.Runtime,
	perms *permission.Resolution,
	in *FormViewInput,
) (view *FormView, err error) {
	cats, err := gh.Repositories(rt).Category.All(ctx, common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	fums, err := gh.Repositories(rt).Forum.All(ctx, common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	view = new(FormView)
	for _, cat := range cats {
		if !perms.CanWriteID(cat.ID) {
			continue
		}
		catfum := CategoryWithForums{
			Category: cat,
		}
		for _, fum := range fums {
			if fum.CategoryID == cat.ID {
				catfum.Forums = append(catfum.Forums, fum)
			}
		}
		view.CategoriesForums = append(view.CategoriesForums, catfum)
	}

	if view.AllowPoll, err = PollAllowed(ctx, rt); err != nil {
		return nil, err
	}
	view.PollOptionsMax = PollOptionsMax

	var writableTopic *vtopic.VTopic
	if in.ForumSlug != "" {
		fum, ferr := gh.Repositories(rt).Forum.VGetBySlug(ctx,
			in.ForumSlug,
			common.QueryOptions{
				Limit: 1,
			},
		)
		if ferr != nil {
			return nil, ferr
		}

		if perms.CanWriteSlug(fum.CategorySlug) {
			view.Forum = fum

			if in.TopicSlug != "" {
				top, terr := gh.Repositories(rt).Topic.VGetByForumUUIDSlug(ctx,
					fum.ID,
					in.TopicSlug,
					common.QueryOptions{
						Limit: 1,
					},
				)
				if terr != nil {
					return nil, terr
				}

				view.Topic = top
				writableTopic = top
			}
		}
	}

	if in.ReplyID != "" && writableTopic != nil {
		rep, rerr := gh.Repositories(rt).Reply.GetByID(ctx,
			in.ReplyID,
			common.QueryOptions{
				Limit: 1,
			},
		)
		if rerr != nil {
			return nil, rerr
		}

		if rep.TopicID == writableTopic.ID {
			view.Reply = rep
		}
	}

	return view, nil
}
