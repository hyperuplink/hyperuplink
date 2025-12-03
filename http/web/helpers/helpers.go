package helpers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GetAbsPathAndRelRoot(c fiber.Ctx) (absPath string, relRoot string) {
	cR := c.Route()
	absPath = cR.Path
	parts := strings.Count(absPath, "/")

	relRoot = ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}

	return absPath, relRoot
}
