package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) notifyReply(
	c fiber.Ctx,
	req *request.Request,
	created *logictopics.CreatedReply,
) {
	cat, err := r.Runtime.Repositories.Category.GetBySlug(
		c.Params("categories"),
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		r.Runtime.Error("reply notification: cannot load category",
			"error", err)
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
		r.Runtime,
		created.Reply,
		req.User.Username,
		req.Ts("reply_notification_subject"),
		loc,
	); err != nil {
		r.Runtime.Error("reply notification: dispatch failed",
			"error", err)
	}
}
