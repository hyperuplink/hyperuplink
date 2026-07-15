package topics

import (
	"fmt"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type ReplyLocation struct {
	CategoryID   uuid.UUID
	CategoryName string
	CategoryPath string

	ForumName string
	ForumPath string

	TopicID   uuid.UUID
	TopicName string
	TopicPath string
}

func (loc ReplyLocation) replyPath(rep *reply.Reply) string {
	return fmt.Sprintf("%s#post-%s", loc.TopicPath, rep.ShortID)
}

func SendReplyNotifications(
	rt *runtime.Runtime,
	rep *reply.Reply,
	byUsername string,
	subject string,
	loc ReplyLocation,
) error {
	recipients, err := rt.Repositories.User.AllToNotifyForReply(
		loc.TopicID,
		rep.ID,
		rep.AuthorID,
		loc.CategoryID,
		common.QueryOptions{},
	)
	if err != nil {
		return err
	}

	if recipients == nil || len(*recipients) == 0 {
		return nil
	}

	payloads := make([]*replynotification.ReplyNotification, 0, len(*recipients))
	for i := range *recipients {
		payload, perr := replynotification.New(&(*recipients)[i], subject)
		if perr != nil {
			return perr
		}

		payload.SetReply(rep, byUsername, loc.replyPath(rep))
		payload.SetCategory(loc.CategoryName, loc.CategoryPath)
		payload.SetForum(loc.ForumName, loc.ForumPath)
		payload.SetTopic(loc.TopicID, loc.TopicName, loc.TopicPath)

		payloads = append(payloads, payload)
	}

	return rt.Dispatch.ReplyNotifications(payloads)
}
