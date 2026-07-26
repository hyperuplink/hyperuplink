package report

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicreport "xn--gckvb8fzb.com/hyperuplink/logic/root/report"
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

	authorID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	err = logicreport.Create(r.Runtime, authorID, &logicreport.CreateInput{
		Target:     frm.Target,
		ID:         frm.ID,
		ReportType: frm.ReportType,
	})
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
