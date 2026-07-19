package newpost

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

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

	in := new(logicnewpost.CreateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", in)

	authorID, _ := req.Session.GetUserUUID()

	var vtop *vtopic.VTopic
	vtop, err = logicnewpost.Create(
		r.Runtime,
		authorID,
		req.Perms(),
		req.Site.GetTimezone(),
		in,
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return helpers.ProcessAttachments(r.Runtime, c, authorID)
		},
	)
	if errors.Is(err, errs.ErrForbidden) {
		return req.RedirectToRoot()
	}
	if err != nil {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute)
	}

	return req.RedirectToRoute(route.For("CategoriesForumsTopics").Fill(
		map[string]string{
			"categories": vtop.CategorySlug,
			"forums":     vtop.ForumSlug,
			"topics":     vtop.Slug,
		},
	))
}
