package report

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type ReportCreateForm struct {
	Target     string `form:"target" validate:"required,oneof=topic reply"`
	ID         string `form:"id" validate:"required"`
	ReportType int    `form:"report_type" validate:"oneof=0 1 2"`
	Return     string `form:"return"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("Report")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ReportCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoot()
	}

	post, err := r.resolvePost(frm.Target, frm.ID)
	if ret, rerr := req.RedirectToRootOnError(err); ret == true {
		return rerr
	}

	authorID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	event := new(postevent.PostEvent)
	event.Type = postevent.Report
	event.AuthorID = authorID
	event.Target = post.Target
	event.TopicID = post.TopicID
	event.ReplyID = post.ReplyID
	event.Selection = frm.ReportType

	_, err = r.Runtime.Repositories.PostEvent.Create(event)
	if errors.Is(err, errs.ErrUniqueViolationOn) {
		req.Flash.SetInfo("report_already")
		return req.RedirectTo(sanitizeReturn(frm.Return))
	}
	if ret, rerr := req.RedirectToRootOnError(err); ret == true {
		return rerr
	}

	return req.RedirectTo(sanitizeReturn(frm.Return))
}

func sanitizeReturn(returnTo string) string {
	if returnTo == "" ||
		strings.Contains(returnTo, "://") ||
		strings.HasPrefix(returnTo, "/") ||
		strings.HasPrefix(returnTo, "\\") {
		return ""
	}

	return returnTo
}
