package request

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request/bcn"
	rerrors "github.com/mrusme/hyperuplink/http/web/request/errors"
	"github.com/mrusme/hyperuplink/http/web/request/form"
	"github.com/mrusme/hyperuplink/http/web/request/in"
	"github.com/mrusme/hyperuplink/http/web/request/session"
	"github.com/mrusme/hyperuplink/http/web/request/site"
)

type Request struct {
	r       route.IRoute
	c       fiber.Ctx
	layouts []string
	view    string
	BCN     *bcn.BreadcrumbNavigation
	Site    *site.Site
	Session *session.Session
	Errors  *rerrors.Errors
	Form    *form.Form
	In      *in.Internationalization
}

func New(r route.IRoute, c fiber.Ctx, layouts []string, view string, title string) *Request {
	req := new(Request)
	req.r = r
	req.c = c
	req.layouts = layouts
	req.view = view
	req.BCN = bcn.New()
	req.Site = site.New(req.r, req.c)
	req.Session = session.New(req.c)
	req.Errors = rerrors.New()
	req.Form = form.New()
	req.In = in.New(req.r, req.c)

	req.Site.SetTitle(req.In.T(title))

	req.BCN.Append(*bcn.NewBreadcrumb(
		true,
		req.Site.Title(),
		req.Site.Title(),
		"",
	))

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
			req.Form.Set(f)
		} else {
			req.Errors.Set(err)
		}

		return false
	}

	return true
}

func (req *Request) RespondWithView(layouts []string, view string) error {
	var layoutsFull []string

	for _, layout := range layouts {
		layoutsFull = append(layoutsFull, fmt.Sprintf("views/layouts/%s", layout))
	}

	return req.c.Render(fmt.Sprintf("views/%s", view), fiber.Map{
		"Breadcrumbs": req.BCN,
		"Site":        req.Site,
		"_":           req.In,
		"Session":     req.Session,
		"Errors":      req.Errors,
		"Form":        req.Form,
	}, layoutsFull...)
}

func (req *Request) Respond() error {
	return req.RespondWithView(req.layouts, req.view)
}

func (req *Request) RespondWithViewOnError(layouts []string, view string, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	req.Errors.Set(err)
	return true, req.RespondWithView(layouts, view)
}

func (req *Request) RedirectTo(path string) error {
	return req.c.Redirect().To(path)
}

func (req *Request) RespondOnError(err error) (bool, error) {
	return req.RespondWithViewOnError(req.layouts, req.view, err)
}

func (req *Request) RedirectToRoot() error {
	return req.RedirectTo(req.Site.GetRelRoot())
}
