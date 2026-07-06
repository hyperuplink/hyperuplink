package manual

import (
	"io/fs"
	"path"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
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

	req.SetData("html", html)

	return req.Respond()
}
