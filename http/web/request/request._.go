package request

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/bcn"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/data"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/flash"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/form"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/in"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/menu"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/session"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type Request struct {
	r       route.IRouteController
	c       fiber.Ctx
	rt      route.Route
	layouts []string
	view    string
	Menu    *menu.Menu
	BCN     *bcn.BreadcrumbNavigation
	Site    *site.Site
	Session *session.Session
	Flash   *flash.Flash
	Form    *form.Form
	In      *in.Internationalization
	Data    *data.Data
	System  *setting.System
	Theme   *setting.Theme
	absPath string
	relRoot string
}

func New(
	r route.IRouteController,
	c fiber.Ctx,
	rt route.Route,
	layouts []string,
	view string,
	title string,
) *Request {
	req := new(Request)
	req.r = r
	req.c = c
	req.rt = rt
	req.layouts = layouts
	req.view = view
	if req.view == "" {
		req.view = "root"
	}
	req.Menu = menu.New()
	req.BCN = bcn.New()
	req.Site = site.New(req.r, req.c)
	req.Session = session.New(req.c)
	req.Flash = flash.New(req.c)
	req.In = in.New(req.r, req.c)
	req.Form = form.New(req.Flash, req.In)
	req.Data = data.New(req.c)

	settingSystem, err := settingRepo.GetByID[setting.System](
		req.r.GetRuntime().Repositories.Setting,
		"system",
	)
	if err != nil {
		req.r.GetRuntime().Error("error", err)
	}
	req.System = &settingSystem.JSONValue

	settingTheme, err := settingRepo.GetByID[setting.Theme](
		req.r.GetRuntime().Repositories.Setting,
		"theme",
	)
	if err != nil {
		req.r.GetRuntime().Error("error", err)
	}
	req.Theme = &settingTheme.JSONValue

	settingGeneral, err := settingRepo.GetByID[setting.General](
		req.r.GetRuntime().Repositories.Setting,
		"general",
	)
	if err != nil {
		req.r.GetRuntime().Error("error", err)
	}

	_, req.absPath, req.relRoot = helpers.GetPaths(req.c)

	if userID, ok := req.Session.GetUserID(); ok {
		usr, err := req.r.GetRuntime().Repositories.User.GetByID(
			userID,
			common.QueryOptions{
				WithBanned:  false,
				WithSpammed: false,
				WithDeleted: false,
				Limit:       1,
			},
		)
		if err != nil {
			// We seemingly have a session but we can't find a user for it in our
			// database. This could be, because maybe the user got banned or deleted.
			// In this case, we destroy the session.
			req.r.GetRuntime().Warn("error", err)
			if err = req.Session.Destroy(); err != nil {
				req.r.GetRuntime().Error("error", err)
			}
		} else {
			req.Session.SetCurrentUser(usr)
		}
	}

	req.Menu.SetI18n(req.In.Ts)
	req.Menu.SetGeneral(settingGeneral.JSONValue)
	req.Menu.SetRole(req.Session.GetCurrentUserRole())

	if title == "" {
		title = req.System.Name
	}
	req.Site.SetRoute(req.rt)
	req.Site.SetTitle(req.In.Ts(title))
	req.Site.SetTimezone(req.Session.GetCurrentUserTimezone())

	var parentRoute route.Route = req.rt
	for parentRoute.Len() > 1 {
		parentRoute = parentRoute.ParentRoute()
		if parentRoute.HasBreadcrumb() {
			req.BCN.Prepend(*bcn.NewBreadcrumb(
				false,
				req.In.Ts(parentRoute.AsTitle()),
				req.In.Ts(parentRoute.AsTitle()),
				req.HrefTo(parentRoute.AsURL()),
			))
		}
	}

	req.BCN.Append(*bcn.NewBreadcrumb(
		true,
		req.Site.Title(),
		req.Site.Title(),
		"",
	))

	return req
}

func (req *Request) UpdateTitle(title string) {
	req.Site.SetTitle(title)
	req.BCN.UpdateLabel(1, title)
}

func (req *Request) UpdateParentTitle(title string) {
	req.BCN.UpdateLabel(2, title)
}

func (req *Request) UpdateParentHref(href string) {
	req.BCN.UpdateHref(2, href)
}

func (req *Request) UpdateGrandParentTitle(title string) {
	req.BCN.UpdateLabel(3, title)
}

func (req *Request) UpdateGrandParentHref(href string) {
	req.BCN.UpdateHref(3, href)
}

func (req *Request) SetData(key string, val interface{}) {
	req.Data.Set(key, val)
}

func (req *Request) HrefTo(path string) string {
	return fmt.Sprintf("%s%s", req.relRoot, path)
}

func (req *Request) ValidateForm(f any, t reflect.Type) bool {
	if err := req.c.Bind().Form(f); err != nil {
		var errmap map[string]error = make(map[string]error)

		if valErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range valErrs {
				if t.Kind() != reflect.Struct {
					req.Flash.SetError(valErrs)
					break
				}

				field, ok := t.FieldByName(e.StructField())
				if !ok {
					req.Flash.SetError(valErrs)
					break
				}

				formTag, ok := field.Tag.Lookup("form")
				if !ok {
					req.Flash.SetError(valErrs)
					break
				}

				errmap[formTag] = fmt.Errorf(
					"%w_%s_%s",
					errs.ErrValidation,
					strings.ToLower(e.Field()),
					e.Tag(),
				)
			}

			req.Flash.SetErrorsMap(errmap)
			req.Form.Set(f)
		} else {
			req.Flash.SetError(err)
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
		"Menu":        req.Menu,
		"Breadcrumbs": req.BCN,
		"Site":        req.Site,
		"_":           req.In,
		"Session":     req.Session,
		"Flash":       req.Flash,
		"Form":        req.Form,
		"Data":        req.Data,
		"System":      req.System,
		"Theme":       req.Theme,
	}, layoutsFull...)
}

func (req *Request) Respond() error {
	return req.RespondWithView(req.layouts, req.view)
}

func (req *Request) RespondWithViewOnError(
	layouts []string,
	view string,
	err error,
) (bool, error) {
	if err == nil {
		return false, nil
	}

	req.Flash.SetError(err)
	return true, req.RespondWithView(layouts, view)
}

func (req *Request) RespondOnError(err error) (bool, error) {
	return req.RespondWithViewOnError(req.layouts, req.view, err)
}

func (req *Request) RespondError(err error) (rerr error) {
	_, rerr = req.RespondWithViewOnError(req.layouts, req.view, err)
	return rerr
}

func (req *Request) redirect(url string) (err error) {
	redir := req.c.Redirect()
	for key, msg := range req.Flash.All() {
		redir = redir.With(key, msg)
	}
	return redir.To(url)
}

func (req *Request) redirectOnError(err error, url string) (bool, error) {
	if err == nil {
		return false, nil
	}

	req.Flash.SetError(err)
	return true, req.redirect(url)
}

func (req *Request) RedirectToRouteID(id string) error {
	return req.redirect(fmt.Sprintf("%s%s",
		req.relRoot,
		route.For(id).AsURL(),
	))
}

func (req *Request) RedirectToRouteIDOnError(err error, id string) (bool, error) {
	return req.redirectOnError(err, fmt.Sprintf("%s%s",
		req.relRoot,
		route.For(id).AsURL(),
	))
}

func (req *Request) RedirectToRoute(r route.Route) error {
	return req.redirect(fmt.Sprintf("%s%s",
		req.relRoot,
		r.AsURL(),
	))
}

func (req *Request) RedirectToRouteOnError(err error, r route.Route) (bool, error) {
	return req.redirectOnError(err, fmt.Sprintf("%s%s",
		req.relRoot,
		r.AsURL(),
	))
}

func (req *Request) RedirectToRouteWithQuery(r route.Route, queries ...string) error {
	url := fmt.Sprintf("%s%s",
		req.relRoot,
		r.AsURL(),
	)
	for i, query := range queries {
		char := "&"
		if i == 0 {
			char = "?"
		} else if i%2 == 0 {
			char = "&"
		} else {
			char = "="
		}
		url = fmt.Sprintf("%s%s%s", url, char, query)
	}
	return req.redirect(url)
}

func (req *Request) RedirectTo(path string) error {
	return req.redirect(fmt.Sprintf("%s%s",
		req.relRoot,
		path,
	))
}

func (req *Request) RedirectToOnError(err error, path string) (bool, error) {
	return req.redirectOnError(err, fmt.Sprintf("%s%s",
		req.relRoot,
		path,
	))
}

func (req *Request) RedirectToRoot() error {
	return req.RedirectTo("")
}

func (req *Request) RedirectToRootOnError(err error) (bool, error) {
	return req.RedirectToOnError(err, "")
}
