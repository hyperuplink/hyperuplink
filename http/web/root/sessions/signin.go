package sessions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/modules/errors"
	"github.com/mrusme/hyperuplink/http/web/modules/session"
	"github.com/mrusme/hyperuplink/http/web/modules/site"
)

type SignInForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignInShow(c fiber.Ctx) error {
	sit := site.New(r, c)

	return c.Render("views/session/signin", fiber.Map{
		"Site": sit,
	}, "views/layouts/base")
}

func (r *Route) SignInCreate(c fiber.Ctx) error {
	sit := site.New(r, c)
	ses := session.New(c)
	ers := errors.New()
	frm := new(SignInForm)

	if errmap, ok := sit.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		ers.SetMap(errmap)

		return c.Render("views/session/signin", fiber.Map{
			"Site":   sit,
			"Errors": ers,
		}, "views/layouts/base")
	}

	r.Runtime.Debug("form", frm)

	if frm.Username == "user" && frm.Password == "pass" {
		if err := ses.Set("local", "2941476f-2ae0-4c3e-a459-1ef5d8dd6ca9"); err != nil {
			ses.Reset()
		}

		return c.Redirect().To(sit.GetRelRoot())
	}

	return c.Render("views/session/signin", fiber.Map{
		"Site": sit,
	}, "views/layouts/base")
}
