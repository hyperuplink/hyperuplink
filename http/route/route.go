package route

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
)

type Route []string

func (r Route) AsURL() string {
	return strings.Join([]string(r), "/")
}

func (r Route) Pathname() string {
	rl := len(r)
	if rl == 0 {
		return ""
	}

	return r[rl-1]
}

var Routes map[string]Route = map[string]Route{
	"Account":                {"account"},
	"AccountSettings":        {"account", "settings"},
	"AccountProfile":         {"account", "profile"},
	"Sessions":               {"sessions"},
	"SessionsConfirm":        {"sessions", "confirm"},
	"SessionsConfirmResend":  {"sessions", "confirm", "resend"},
	"SessionsSignup":         {"sessions", "signun"},
	"SessionsSignin":         {"sessions", "signin"},
	"SessionsSignout":        {"sessions", "signout"},
	"Admin":                  {"admin"},
	"Categories":             {"categories"},
	"CategoriesForums":       {"categories", "forums"},
	"CategoriesForumsTopics": {"categories", "forums", "topics"},
	"System":                 {"system"},
}

type IRouteController interface {
	GetRuntime() *runtime.Runtime
	GetPath() string
	GetEnv() *Environment

	// Index(fiber.Ctx) error
	// Show(fiber.Ctx) error
	// Create(fiber.Ctx) error
	// Update(fiber.Ctx) error
	// Destroy(fiber.Ctx) error
}

type Environment struct {
	Title string
}

type RouteController struct {
	Runtime *runtime.Runtime
	Router  fiber.Router
	Routes  []IRouteController
	Path    string
	Env     *Environment
}

func NewEnv() *Environment {
	env := new(Environment)

	env.Title = "Hyperuplink"

	return env
}

func For(id string) (r Route) {
	var ok bool

	if r, ok = Routes[id]; !ok {
		return Route{}
	}

	return r
}

func GetReservedBasePaths(app *fiber.App) []string {
	var reserved []string

	froutes := app.GetRoutes(true)
	for _, r := range froutes {
		sr := strings.Split(r.Path, "/")
		if len(sr) >= 2 {
			if sr[1][0] != ':' {
				reserved = append(reserved, sr[1])
			}
		}
	}

	return reserved
}
