package sessions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	// "github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	s := site.New(r, c)

	return c.Render("views/session/signup", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	s := site.New(r, c)
	// sess := session.FromContext(c)
	f := new(SignUpForm)

	if err := c.Bind().Form(f); err != nil {
		var errs map[string]error = make(map[string]error)

		if valErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range valErrs {
				t := reflect.TypeOf(*f)
				if t.Kind() != reflect.Struct {
					// TODO: errrrrr.....?
				}

				field, ok := t.FieldByName(e.StructField())
				if !ok {
					// TODO: hmpf
				}

				formTag, ok := field.Tag.Lookup("form")
				if !ok {
					// TODO: hmpffff
				}

				errs[formTag] = errors.New(s.T(fmt.Sprintf(
					"validation_%s_%s",
					strings.ToLower(e.Field()),
					e.Tag(),
				)))
			}
		} else {
			errs["error"] = err
		}

		return c.Render("views/session/signup", fiber.Map{
			"Site":   s,
			"Errors": errs,
		}, "views/layouts/base")
	}

	r.Runtime.Debug("form", f)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return c.Render("views/session/signup", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}
