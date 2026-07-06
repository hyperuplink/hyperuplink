package topics

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type TopicsUpdateForm struct {
	AllowKindQuestion bool `form:"allow_kind_question"`
	AllowKindPoll     bool `form:"allow_kind_poll"`
	AllowKindRSVP     bool `form:"allow_kind_rsvp"`
}

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

	frm := new(TopicsUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	var settingTopics *setting.Setting[setting.Topics]
	settingTopics, err = settingRepo.GetByID[setting.Topics](
		r.Runtime.Repositories.Setting,
		"topics",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingTopics.JSONValue.AllowKindQuestion = frm.AllowKindQuestion
	settingTopics.JSONValue.AllowKindPoll = frm.AllowKindPoll
	settingTopics.JSONValue.AllowKindRSVP = frm.AllowKindRSVP

	err = settingRepo.Update[setting.Topics](
		r.Runtime.Repositories.Setting,
		settingTopics,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
