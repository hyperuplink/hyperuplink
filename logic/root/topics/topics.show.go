package topics

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type ShowInput struct {
	ID       string
	ViewerID uuid.NullUUID
}

type Shown struct {
	Topic  *vtopic.VTopic
	Poll   *logictopics.Poll
	Unread bool
}

func Resolve(
	rt *runtime.Runtime,
	idOrShortID string,
) (top *vtopic.VTopic, err error) {
	if id, perr := uuid.Parse(idOrShortID); perr == nil {
		return gh.Repositories(rt).Topic.VGetByUUID(id, common.QueryOptions{Limit: 1})
	}

	return gh.Repositories(rt).Topic.VGetByShortID(idOrShortID, common.QueryOptions{Limit: 1})
}

func Read(
	rt *runtime.Runtime,
	idOrShortID string,
	perms *permission.Resolution,
) (top *vtopic.VTopic, err error) {
	if top, err = Resolve(rt, idOrShortID); err != nil {
		return nil, err
	}

	if !perms.CanReadSlug(top.CategorySlug) {
		return nil, errs.ErrForbidden
	}

	return top, nil
}

func Show(
	rt *runtime.Runtime,
	in *ShowInput,
	perms *permission.Resolution,
) (shown *Shown, err error) {
	top, err := Read(rt, in.ID, perms)
	if err != nil {
		return nil, err
	}

	shown = &Shown{Topic: top}

	if in.ViewerID.Valid {
		var unread map[string]bool
		if unread, err = logicactivity.UnreadTopics(rt, in.ViewerID.UUID, &[]vtopic.VTopic{*top}); err != nil {
			return nil, err
		}
		shown.Unread = unread[top.ID.String()]

		logicactivity.RecordTopicView(rt, in.ViewerID.UUID, top.ID)
	}

	if shown.Poll, err = logictopics.PollView(rt, &logictopics.PollViewInput{
		Topic:    &top.Topic,
		ViewerID: in.ViewerID,
		CanWrite: perms.CanWriteSlug(top.CategorySlug),
	}); err != nil {
		return nil, err
	}

	return shown, nil
}
