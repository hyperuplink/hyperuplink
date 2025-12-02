package request

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	rerrors "github.com/mrusme/hyperuplink/http/web/request/errors"
	"github.com/mrusme/hyperuplink/http/web/request/session"
	"github.com/mrusme/hyperuplink/http/web/request/site"
)

type Request struct {
	r       route.IRoute
	c       fiber.Ctx
	Site    *site.Site
	Session *session.Session
	Errors  *rerrors.Errors
}

func New(r route.IRoute, c fiber.Ctx) *Request {
	req := new(Request)
	req.r = r
	req.c = c
	req.Site = site.New(req.r, req.c)
	req.Session = session.New(req.c)
	req.Errors = rerrors.New()

	return req
}

func (req *Request) ValidateForm(f any, t reflect.Type) bool {
	if err := req.c.Bind().Form(f); err != nil {
		var errmap map[string]error = make(map[string]error)

		if valErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range valErrs {
				if t.Kind() != reflect.Struct {
					req.Errors.Set(valErrs)
					break
				}

				field, ok := t.FieldByName(e.StructField())
				if !ok {
					req.Errors.Set(valErrs)
					break
				}

				formTag, ok := field.Tag.Lookup("form")
				if !ok {
					req.Errors.Set(valErrs)
					break
				}

				errmap[formTag] = errors.New(fmt.Sprintf(
					"validation_%s_%s",
					strings.ToLower(e.Field()),
					e.Tag(),
				))
			}

			req.Errors.SetMap(errmap)
		} else {
			req.Errors.Set(err)
		}

		return false
	}

	return true
}

func (req *Request) Respond(layout string, view string) error {
	return req.c.Render(fmt.Sprintf("views/%s", view), fiber.Map{
		"Site":    req.Site,
		"Session": req.Session,
		"Errors":  req.Errors,
	}, fmt.Sprintf("views/layouts/%s", layout))
}

func (req *Request) RedirectTo(path string) error {
	return req.c.Redirect().To(path)
}

func (req *Request) RedirectToRoot() error {
	return req.RedirectTo(req.Site.GetRelRoot())
}
