package manual

import (
	"io/fs"
	"path"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/http/web/request/bcn"
	"github.com/mrusme/hyperuplink/models/user"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	return r.render(c, c.Params("*"))
}

func (r *Route) render(c fiber.Ctx, sub string) (err error) {
	embedPath := path.Join("docs", "manual", strings.Trim(sub, "/"), "index.md")
	if embedPath != "docs/manual/index.md" &&
		!strings.HasPrefix(embedPath, "docs/manual/") {
		return c.SendStatus(fiber.StatusNotFound)
	}

	reqPath := c.Path()
	if !strings.HasSuffix(reqPath, "/") {
		return c.Redirect().
			Status(fiber.StatusMovedPermanently).
			To(reqPath[strings.LastIndexByte(reqPath, '/')+1:] + "/")
	}

	myRoute := route.For("DocsManual")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.ModRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var src []byte
	src, err = fs.ReadFile(r.Runtime.Embeds["docs"], embedPath)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var html string
	html, err = r.Runtime.Markdown.Convert(string(src))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	r.setSegmentBreadcrumbs(req, myRoute, sub)

	req.SetData("html", html)

	return req.Respond()
}

func (r *Route) setSegmentBreadcrumbs(
	req *request.Request,
	myRoute route.Route,
	sub string,
) {
	var segments []string
	for _, segment := range strings.Split(strings.Trim(sub, "/"), "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	if len(segments) == 0 {
		return
	}

	crumbs := req.BCN.Get()
	if n := len(crumbs); n > 0 {
		crumbs[n-1].IsActive = false
		crumbs[n-1].Href = req.HrefTo(myRoute.AsURL() + "/")
	}

	acc := myRoute.AsURL()
	for i, segment := range segments {
		acc += "/" + segment

		href := ""
		if i < len(segments)-1 {
			href = req.HrefTo(acc + "/")
		}

		label := strings.ToUpper(segment[:1]) + segment[1:]
		crumbs = append(crumbs, *bcn.NewBreadcrumb(
			i == len(segments)-1,
			label,
			label,
			href,
		))
	}

	req.BCN.Set(crumbs)
}
