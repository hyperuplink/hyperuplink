package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/runtime"
)

type IRoute interface {
	Index(fiber.Ctx) error
	Show(fiber.Ctx) error
	Create(fiber.Ctx) error
	Update(fiber.Ctx) error
	Destroy(fiber.Ctx) error
}

type Route struct {
	Runtime *runtime.Runtime
	Router  fiber.Router
	Routes  []IRoute
}
