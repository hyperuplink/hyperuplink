package newpost

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type NewCreateForm struct {
	Name        string   `form:"name" validate:"required,min=1,max=78"`
	Text        string   `form:"text" validate:"required,min=1"`
	ForumID     string   `form:"forum_id" validate:"required,uuid"`
	Kind        string   `form:"kind" validate:"omitempty,oneof=regular poll"`
	PollOptions []string `form:"poll_options" validate:"omitempty,dive,max=78"`
	PollEndsAt  string   `form:"poll_ends_at"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
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

	frm := new(NewCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	top := new(topic.Topic)
	top.ShortID = shortuuid.New()
	top.Name = frm.Name
	top.SetSlugFromName()
	top.ForumID, err = uuid.Parse(frm.ForumID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	fum, err := r.Runtime.Repositories.Forum.GetByUUID(
		top.ForumID,
		common.QueryOptions{Limit: 1},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if !req.Perms().CanWriteID(fum.CategoryID) {
		return req.RedirectToRoot()
	}

	top.AuthorID, _ = req.Session.GetUserUUID()
	top.Kind = topic.Regular
	top.Anonymous = false
	top.Text = frm.Text
	top.HTML, err = r.Runtime.Markdown.Convert(top.Text)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if topic.Kind(frm.Kind) == topic.Poll {
		err = logicnewpost.ApplyPoll(r.Runtime, top, &logicnewpost.PollInput{
			Options:  frm.PollOptions,
			EndsAt:   frm.PollEndsAt,
			Location: req.Site.GetTimezone(),
		})
		if err != nil {
			req.Flash.SetError(err)
			return req.RedirectToRoute(myRoute)
		}
	}

	top.AttachmentIDs, err = helpers.ProcessAttachments(r.Runtime, c, top.AuthorID)
	if err != nil {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute)
	}

	_, err = r.Runtime.Repositories.Topic.Create(top)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var vtop *vtopic.VTopic
	vtop, err = r.Runtime.Repositories.Topic.VGetByForumUUIDSlug(
		top.ForumID,
		top.Slug,
		common.QueryOptions{
			Limit: 1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(route.For("CategoriesForumsTopics").Fill(
		map[string]string{
			"categories": vtop.CategorySlug,
			"forums":     vtop.ForumSlug,
			"topics":     top.Slug,
		},
	))
}
