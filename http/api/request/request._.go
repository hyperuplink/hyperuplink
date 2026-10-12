package request

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/intnat"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicperms "xn--gckvb8fzb.com/hyperuplink/logic/helpers/perms"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const (
	PrincipalLocal = "api_principal"
	SystemLocal    = "api_system"

	MaxPerPage = 100
)

type Principal struct {
	User *user.User
	Key  *apikey.APIKey
}

func (p Principal) IsGuest() bool {
	return p.User == nil
}

func (p Principal) Role() user.Role {
	if p.User == nil {
		return user.GuestRole
	}
	return p.User.Role
}

func PrincipalFrom(c fiber.Ctx) Principal {
	if p, ok := c.Locals(PrincipalLocal).(Principal); ok {
		return p
	}
	return Principal{}
}

func SystemFrom(c fiber.Ctx, rt *runtime.Runtime) *setting.System {
	if sys, ok := c.Locals(SystemLocal).(*setting.System); ok && sys != nil {
		return sys
	}

	settingSystem, err := settingRepo.GetByID[setting.System](c.Context(),
		gh.Repositories(rt).Setting,
		"system",
	)
	if err != nil {
		rt.Error("error", err)
		sys := setting.NewSystem()
		return &sys
	}

	sys := settingSystem.JSONValue
	c.Locals(SystemLocal, &sys)

	return &sys
}

type Request struct {
	r         route.IRouteController
	c         fiber.Ctx
	Principal Principal
	System    *setting.System
	perms     *permission.Resolution
	views     *view.Mapper
	lang      string
}

func New(
	r route.IRouteController,
	c fiber.Ctx,
) *Request {
	req := new(Request)
	req.r = r
	req.c = c
	req.Principal = PrincipalFrom(c)
	req.System = SystemFrom(c, r.GetRuntime())

	al := c.Get("Accept-Language", "en")
	als := strings.Split(al, "-")
	req.lang = strings.ToLower(als[0])

	return req
}

func (req *Request) Ctx() fiber.Ctx {
	return req.c
}

func (req *Request) Ts(msg string) string {
	return intnat.From(req.r.GetRuntime()).NewLocalizer(req.lang).Get(msg)
}

func (req *Request) Role() user.Role {
	return req.Principal.Role()
}

func (req *Request) IsGuest() bool {
	return req.Principal.IsGuest()
}

func (req *Request) User() *user.User {
	return req.Principal.User
}

func (req *Request) UserUUID() (uuid.UUID, bool) {
	if req.Principal.User == nil {
		return uuid.Nil, false
	}
	return req.Principal.User.ID, true
}

func (req *Request) ViewerID() (viewerID uuid.NullUUID) {
	viewerID.UUID, viewerID.Valid = req.UserUUID()
	return viewerID
}

func (req *Request) Perms() *permission.Resolution {
	if req.perms != nil {
		return req.perms
	}

	var memberOf []string
	if req.Principal.User != nil {
		memberOf = req.Principal.User.MemberOf
	}

	req.perms = logicperms.Resolve(
		req.Context(),
		req.r.GetRuntime(),
		req.Role(),
		memberOf,
	)
	return req.perms
}

func (req *Request) Views() *view.Mapper {
	if req.views == nil {
		req.views = view.NewMapper(req.Context(), req.r.GetRuntime(), req.System)
	}
	return req.views
}

func (req *Request) RequireUser() (mustReturn bool, err error) {
	if req.IsGuest() {
		return true, req.RespondError(errs.ErrAuthenticationRequired)
	}

	return false, nil
}

func (req *Request) RequireAdmin() (mustReturn bool, err error) {
	if req.IsGuest() {
		return true, req.RespondError(errs.ErrAuthenticationRequired)
	}
	if req.Role() != user.AdminRole {
		return true, req.RespondError(errs.ErrForbidden)
	}

	return false, nil
}

func (req *Request) Page(defaultPerPage int) (page int, perPage int, mustReturn bool, err error) {
	fields := map[string][]string{}

	page = 1
	if raw := req.c.Query("page"); raw != "" {
		if page, err = strconv.Atoi(raw); err != nil {
			fields["page"] = []string{"number"}
		} else if page < 1 {
			fields["page"] = []string{"min"}
		}
	}

	perPage = defaultPerPage
	if raw := req.c.Query("per_page"); raw != "" {
		if perPage, err = strconv.Atoi(raw); err != nil {
			fields["per_page"] = []string{"number"}
		} else if perPage < 1 {
			fields["per_page"] = []string{"min"}
		} else if perPage > MaxPerPage {
			fields["per_page"] = []string{"max"}
		}
	}

	if len(fields) > 0 {
		return 0, 0, true, req.RespondProblem(problem.Validation(fields))
	}

	return page, perPage, false, nil
}

func (req *Request) Sort(fallback string, allowed ...string) (sort string, mustReturn bool, err error) {
	sort = req.c.Query("sort", fallback)
	if !slices.Contains(allowed, sort) {
		return "", true, req.RespondProblem(problem.Validation(map[string][]string{
			"sort": {"oneof"},
		}))
	}

	return sort, false, nil
}

func (req *Request) ParamUUID(name string) (id uuid.UUID, mustReturn bool, err error) {
	if id, err = uuid.Parse(req.c.Params(name)); err != nil {
		return uuid.Nil, true, req.RespondError(errs.ErrNoRows)
	}

	return id, false, nil
}

func (req *Request) BindJSON(in any) (mustReturn bool, err error) {
	if berr := req.c.Bind().JSON(in); berr != nil {
		return true, req.respondBindError(in, berr)
	}

	return false, nil
}

func (req *Request) respondBindError(in any, berr error) error {
	var valErrs validator.ValidationErrors
	if errors.As(berr, &valErrs) {
		fields := map[string][]string{}
		for _, e := range valErrs {
			name := jsonFieldName(reflect.TypeOf(in), e.StructNamespace())
			fields[name] = append(fields[name], e.Tag())
		}

		return req.RespondProblem(problem.Validation(fields))
	}

	return req.RespondError(errs.ErrFormInvalid)
}

func jsonFieldName(t reflect.Type, namespace string) string {
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}

	var field reflect.StructField
	for _, part := range parts {
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return strings.ToLower(parts[len(parts)-1])
		}

		var ok bool
		if field, ok = t.FieldByName(part); !ok {
			return strings.ToLower(parts[len(parts)-1])
		}
		t = field.Type
	}

	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" || name == "-" {
		return strings.ToLower(field.Name)
	}

	return name
}

func (req *Request) Respond(data any) error {
	return req.c.JSON(data)
}

func (req *Request) RespondCreated(data any) error {
	return req.c.Status(fiber.StatusCreated).JSON(data)
}

func (req *Request) RespondNoContent() error {
	return req.c.SendStatus(fiber.StatusNoContent)
}

func (req *Request) RespondProblem(p problem.Problem) error {
	return problem.Respond(req.c, p)
}

func (req *Request) RespondError(err error) error {
	return req.RespondProblem(problem.From(err))
}

func (req *Request) RespondOnReadError(err error) (mustReturn bool, rerr error) {
	if errors.Is(err, errs.ErrForbidden) {
		err = errs.ErrNoRows
	}

	return req.RespondOnError(err)
}

func (req *Request) RespondOnError(err error) (mustReturn bool, rerr error) {
	if err == nil {
		return false, nil
	}

	return true, req.RespondError(err)
}

func (req *Request) Context() context.Context {
	return req.c.Context()
}
