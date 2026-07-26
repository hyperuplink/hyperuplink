package topics

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardTopics")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logictopics.UpdateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", in)

	actorID := uuid.NullUUID{}
	actorID.UUID, actorID.Valid = req.Session.GetUserUUID()

	err = logictopics.Update(r.Runtime, actorID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
