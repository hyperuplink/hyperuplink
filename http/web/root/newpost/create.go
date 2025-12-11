package newpost

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/topic"
	"github.com/mrusme/hyperuplink/models/user"
)

type NewCreateForm struct {
	Name    string `form:"name" validate:"required,min=1,max=78"`
	Text    string `form:"text" validate:"required,min=1"`
	ForumID string `form:"forum_id" validate:"required,uuid"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("New")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
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
	top.Name = frm.Name
	top.SetSlugFromName()
	top.ForumID, err = uuid.Parse(frm.ForumID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	top.AuthorID, _ = req.Session.GetUserUUID()
	top.Kind = topic.Regular
	top.Anonymous = false
	top.Text = frm.Text
	// top.PollOptions =
	_, err = r.Runtime.Repositories.Topic.Create(top)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	fill := make(map[string]string)
	fill["categories"] = "category" // TODO: Set category slug
	fill["forums"] = "forum"       // TODO: Set forum slug
	fill["topics"] = top.Slug
	return req.RedirectToRoute(route.For("CategoriesForumsTopics").Fill(fill))
}
