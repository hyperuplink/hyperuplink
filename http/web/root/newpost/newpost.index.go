package newpost

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("New")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var category_slug string = c.Query("category")
	var forum_slug string = c.Query("forum")

	if category_slug != "" && forum_slug != "" {
		req.SetData("select_category_slug", category_slug)
		req.SetData("select_forum_slug", forum_slug)
	}

	var view *logicnewpost.FormView
	view, err = logicnewpost.View(r.Runtime, req.Perms(), &logicnewpost.FormViewInput{
		ForumSlug: forum_slug,
		TopicSlug: c.Query("topic"),
		ReplyID:   c.Query("reply"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("categories_forums", view.CategoriesForums)
	req.SetData("allow_poll", view.AllowPoll)
	req.SetData("poll_options_max", view.PollOptionsMax)

	if view.Forum != nil {
		req.SetData("forum", view.Forum)
	}
	if view.Topic != nil {
		req.SetData("topic", view.Topic)
	}
	if view.Reply != nil {
		req.SetData("reply", view.Reply)
	}

	return req.Respond()
}
