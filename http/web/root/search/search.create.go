package search

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type SearchForm struct {
	Query        string `form:"q" validate:"required,min=1,max=256"`
	Author       string `form:"author" validate:"omitempty,max=32"`
	FTitle       bool   `form:"f_title"`
	FBody        bool   `form:"f_body"`
	FReplies     bool   `form:"f_replies"`
	FAttachments bool   `form:"f_attachments"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("Search")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Only if forum is public
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(SearchForm)
	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	params := url.Values{}
	params.Set("q", frm.Query)
	if author := strings.TrimSpace(frm.Author); author != "" {
		params.Set("author", author)
	}
	if frm.FTitle {
		params.Set("f_title", "1")
	}
	if frm.FBody {
		params.Set("f_body", "1")
	}
	if frm.FReplies {
		params.Set("f_replies", "1")
	}
	if frm.FAttachments {
		params.Set("f_attachments", "1")
	}

	return req.RedirectTo("search?" + params.Encode())
}
