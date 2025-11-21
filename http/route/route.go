package route

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
)

var (
	AccountRoute                []string = []string{"account"}
	AccountSettingsRoute        []string = []string{"account", "settings"}
	AccountProfileRoute         []string = []string{"account", "profile"}
	SessionRoute                []string = []string{"session"}
	AdminRoute                  []string = []string{"admin"}
	CategoriesRoute             []string = []string{"categories"}
	CategoriesForumsRoute       []string = []string{"categories", "forums"}
	CategoriesForumsTopicsRoute []string = []string{"categories", "forums", "topics"}
	SystemRoute                 []string = []string{"system"}
)

type IRoute interface {
	GetRuntime() *runtime.Runtime
	GetPath() string
	GetEnv() *Environment

	Index(fiber.Ctx) error
	Show(fiber.Ctx) error
	Create(fiber.Ctx) error
	Update(fiber.Ctx) error
	Destroy(fiber.Ctx) error
}

type Environment struct {
	Title string
}

type Route struct {
	Runtime *runtime.Runtime
	Router  fiber.Router
	Routes  []IRoute
	Path    string
	Env     *Environment
}

func NewEnv() *Environment {
	env := new(Environment)

	env.Title = "Hyperuplink"

	return env
}

func GetPathOf(r []string) string {
	return r[len(r)-1]
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
