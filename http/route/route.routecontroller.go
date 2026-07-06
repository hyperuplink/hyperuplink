package route

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

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

type RouteController struct {
	Runtime *runtime.Runtime
	Router  fiber.Router
	Routes  []IRouteController
	Path    string
	Env     *Environment
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
