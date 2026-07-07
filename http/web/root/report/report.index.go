package report

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
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

	target := c.Query("target")
	id := c.Query("id")
	returnTo := c.Query("return")

	post, err := r.resolvePost(target, id)
	if ret, rerr := req.RedirectToRootOnError(err); ret == true {
		return rerr
	}

	req.SetData("target", target)
	req.SetData("id", id)
	req.SetData("return", returnTo)
	req.SetData("post_text", post.Text)

	return req.Respond()
}
