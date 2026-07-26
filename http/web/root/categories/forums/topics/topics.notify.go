package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

func (r *Route) notifyReply(
	c fiber.Ctx,
	req *request.Request,
	rep *reply.Reply,
	top *topic.Topic,
	fum *forum.Forum,
) {
	cat, err := gh.Repositories(r.Runtime).Category.GetBySlug(
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

		ForumName: fum.Name,
		ForumPath: route.For("CategoriesForums").Fill(map[string]string{
			"categories": cat.Slug,
			"forums":     fum.Slug,
		}).AsURL(),

		TopicID:   top.ID,
		TopicName: top.Name,
		TopicPath: route.For("CategoriesForumsTopics").Fill(map[string]string{
			"categories": cat.Slug,
			"forums":     fum.Slug,
			"topics":     top.Slug,
		}).AsURL(),
	}

	if err = logictopics.SendReplyNotifications(
		r.Runtime,
		rep,
		req.Session.GetCurrentUserUsername(),
		req.In.Ts("reply_notification_subject"),
		loc,
	); err != nil {
		r.Runtime.Error("reply notification: dispatch failed",
			"error", err)
	}
}
