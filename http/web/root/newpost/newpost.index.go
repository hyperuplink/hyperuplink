package newpost

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type CategoryWithForums struct {
	Category category.Category
	Forums   []forum.Forum
}

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
	var topic_slug string = c.Query("topic")
	var reply_id string = c.Query("reply")

	if category_slug != "" && forum_slug != "" {
		req.SetData("select_category_slug", category_slug)
		req.SetData("select_forum_slug", forum_slug)
	}

	var cats *[]category.Category
	cats, err = r.Runtime.Repositories.Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var fums *[]forum.Forum
	fums, err = r.Runtime.Repositories.Forum.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var catsfums []CategoryWithForums
	for _, cat := range *cats {
		catfum := CategoryWithForums{
			Category: cat,
		}
		for _, fum := range *fums {
			if fum.CategoryID == cat.ID {
				catfum.Forums = append(catfum.Forums, fum)
			}
		}
		catsfums = append(catsfums, catfum)
	}

	req.SetData("categories_forums", catsfums)

	if forum_slug != "" {
		var fum *vforum.VForum
		fum, err = r.Runtime.Repositories.Forum.VGetBySlug(
			forum_slug,
			common.QueryOptions{
				Limit: 1,
			},
		)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}

		req.SetData("forum", fum)

		if topic_slug != "" {
			var top *vtopic.VTopic
			top, err = r.Runtime.Repositories.Topic.VGetByForumUUIDSlug(
				fum.ID,
				topic_slug,
				common.QueryOptions{
					Limit: 1,
				},
			)
			if ret, rerr := req.RespondOnError(err); ret == true {
				return rerr
			}

			req.SetData("topic", top)
		}
	}

	if reply_id != "" {

		var rep *reply.Reply
		rep, err = r.Runtime.Repositories.Reply.GetByID(
			reply_id,
			common.QueryOptions{
				Limit: 1,
			},
		)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}

		req.SetData("reply", rep)
	}

	// req.UpdateTitle(fum.Name)
	// req.UpdateParentHref(req.HrefTo("_" + fum.CategorySlug))
	// req.UpdateParentTitle(fum.CategoryName)

	return req.Respond()
}
