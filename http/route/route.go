package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
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
