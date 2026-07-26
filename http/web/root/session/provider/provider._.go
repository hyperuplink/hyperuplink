package provider

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/markbates/goth"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type Route struct {
	route.RouteController
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("SessionProvider").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("",
			goth_fiber.BeginAuthHandler).Name("provider.show")
		base.Get("/callback",
			r.ProviderCallbackShow).Name("provider.callback.show")
	}, r.Path+".")

	return r, nil
}

func (r *Route) GetRuntime() *runtime.Runtime {
	return r.Runtime
}

func (r *Route) GetPath() string {
	return r.Path
}

func (r *Route) GetEnv() *route.Environment {
	return r.Env
}

func (r *Route) ProviderCallbackShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionProviderCallback")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	gothUser, err := goth_fiber.CompleteUserAuth(c,
		goth_fiber.CompleteUserAuthOptions{ShouldLogout: false},
	)
	if err != nil {
		r.Runtime.Warn("oauth", "complete auth failed", "error", err)
		req.Flash.SetError(errs.ErrOAuthFailed)
		return req.RedirectToRouteID("SessionSignin")
	}

	if gothUser.Email == "" {
		req.Flash.SetError(errs.ErrOAuthNoEmail)
		return req.RedirectToRouteID("SessionSignin")
	}

	if !oauthEmailVerified(gothUser) {
		r.Runtime.Warn("oauth", "unverified provider email",
			"provider", gothUser.Provider, "email", gothUser.Email)
		req.Flash.SetError(errs.ErrOAuthEmailUnverified)
		return req.RedirectToRouteID("SessionSignin")
	}

	var usr *user.User
	usr, err = gh.Repositories(r.Runtime).User.GetByEmail(
		gothUser.Email,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		if !errors.Is(err, errs.ErrNoRows) {
			req.Flash.SetError(err)
			return req.RedirectToRouteID("SessionSignin")
		}

		usr, err = r.createOAuthUser(req, gothUser)
		if err != nil {
			r.Runtime.Error("oauth", "user creation failed", "error", err)
			req.Flash.SetError(errs.ErrOAuthFailed)
			return req.RedirectToRouteID("SessionSignin")
		}
	}

	if usr.OTPEnabled {
		req.Session.SetPending2FA(gothUser.Provider, usr.ID.String())
		return req.RedirectToRouteID("SessionTwofactor")
	}

	if err = req.Session.Set(gothUser.Provider, usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	return req.RedirectToRoot()
}

func (r *Route) createOAuthUser(
	req *request.Request,
	gothUser goth.User,
) (*user.User, error) {
	usr := new(user.User)
	usr.Username = r.uniqueUsername(gothUser)
	usr.Role = user.UserRole
	usr.Email = gothUser.Email
	usr.Language = req.In.Lang()

	if err := usr.SetRandomPassword(); err != nil {
		return nil, err
	}

	id, err := gh.Repositories(r.Runtime).User.Create(usr)
	if err != nil {
		return nil, err
	}

	usr, err = gh.Repositories(r.Runtime).User.GetByUUID(
		id,
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	usr.SetConfirmedAt(now)
	usr.SetEmailConfirmedAt(now)
	if err := gh.Repositories(r.Runtime).User.Update(usr); err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *Route) uniqueUsername(gothUser goth.User) string {
	base := sanitizeUsername(gothUser.NickName)
	if base == "" {
		base = sanitizeUsername(gothUser.Name)
	}
	if base == "" {
		if at := strings.IndexByte(gothUser.Email, '@'); at > 0 {
			base = sanitizeUsername(gothUser.Email[:at])
		}
	}
	if len(base) < 2 {
		base = "user"
	}

	candidate := base
	for i := 1; i <= 9999; i++ {
		if _, err := gh.Repositories(r.Runtime).User.GetByUsername(
			candidate,
			common.QueryOptions{Limit: 1},
		); err != nil {
			return candidate
		}

		suffix := strconv.Itoa(i)
		trimmed := base
		if len(trimmed)+len(suffix) > 32 {
			trimmed = trimmed[:32-len(suffix)]
		}
		candidate = trimmed + suffix
	}

	return candidate
}

func oauthEmailVerified(gothUser goth.User) bool {
	switch strings.ToLower(gothUser.Provider) {
	case "github", "facebook":
		return true
	case "google":
		return rawDataFlag(gothUser.RawData, "email_verified") ||
			rawDataFlag(gothUser.RawData, "verified_email")
	default:
		return false
	}
}

func rawDataFlag(raw map[string]interface{}, key string) bool {
	value, ok := raw[key]
	if !ok {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	}

	return false
}

func sanitizeUsername(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			b.WriteRune(r)
		}
	}

	out := b.String()
	if len(out) > 32 {
		out = out[:32]
	}

	return out
}
