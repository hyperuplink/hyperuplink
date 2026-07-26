package user

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type View struct {
	User    *user.User
	Groups  *[]group.Group
	Topics  []vtopic.VTopic
	Replies []vreply.VReplyWithTopic
}

func Show(
	rt *runtime.Runtime,
	username string,
	perms *permission.Resolution,
) (view *View, err error) {
	usr, err := gh.Repositories(rt).User.GetByUsername(
		username,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		return nil, err
	}

	groups, err := gh.Repositories(rt).Group.All(common.QueryOptions{
		OrderBy: "name",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	topics, err := gh.Repositories(rt).Topic.VAllForAuthorUUID(
		usr.ID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
			Limit:   10,
		},
	)
	if err != nil {
		return nil, err
	}

	visibleTopics := []vtopic.VTopic{}
	for _, top := range *topics {
		if perms.CanReadSlug(top.CategorySlug) {
			visibleTopics = append(visibleTopics, top)
		}
	}

	replies, err := gh.Repositories(rt).Reply.VAllWithTopicForAuthorUUID(
		usr.ID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Descending,
			Limit:   10,
		},
	)
	if err != nil {
		return nil, err
	}

	visibleReplies := []vreply.VReplyWithTopic{}
	for _, rep := range *replies {
		if perms.CanReadSlug(rep.CategorySlug) {
			visibleReplies = append(visibleReplies, rep)
		}
	}

	return &View{
		User:    usr,
		Groups:  groups,
		Topics:  visibleTopics,
		Replies: visibleReplies,
	}, nil
}
