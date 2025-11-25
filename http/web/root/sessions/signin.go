package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignInForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (r *Route) SignInShow(c fiber.Ctx) error {
	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}

func (r *Route) SignInCreate(c fiber.Ctx) error {
	s := site.New(r, c)
	sess := session.FromContext(c)
	f := new(SignInForm)

	if err := c.Bind().Form(f); err != nil {
		return err // TODO
	}

	r.Runtime.Debug("form", f)
	if f.Username == "user" && f.Password == "pass" {
		if err := sess.Regenerate(); err != nil {
			return err // TODO
		}

		sess.Set("user_id", "2941476f-2ae0-4c3e-a459-1ef5d8dd6ca9")
		sess.Set("auth", "local")

		return c.Redirect().To(s.GetRelRoot())
	}

	return c.Render("views/session/signin", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}
