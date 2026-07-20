package notify

import (
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func Reply(
	rt *runtime.Runtime,
	created *logictopics.CreatedReply,
	byUsername string,
	subject string,
) {
	cat, err := rt.Repositories.Category.GetByUUID(
		created.Forum.CategoryID,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		rt.Error("reply notification: cannot load category", "error", err)
		return
	}

	loc := logictopics.ReplyLocation{
		CategoryID:   cat.ID,
		CategoryName: cat.Name,
		CategoryPath: route.For("Categories").Fill(map[string]string{
			"categories": cat.Slug,
		}).AsURL(),

		ForumName: created.Forum.Name,
		ForumPath: route.For("CategoriesForums").Fill(map[string]string{
			"categories": cat.Slug,
			"forums":     created.Forum.Slug,
		}).AsURL(),

		TopicID:   created.Topic.ID,
		TopicName: created.Topic.Name,
		TopicPath: route.For("CategoriesForumsTopics").Fill(map[string]string{
			"categories": cat.Slug,
			"forums":     created.Forum.Slug,
			"topics":     created.Topic.Slug,
		}).AsURL(),
	}

	if err = logictopics.SendReplyNotifications(
		rt,
		created.Reply,
		byUsername,
		subject,
		loc,
	); err != nil {
		rt.Error("reply notification: dispatch failed", "error", err)
	}
}
