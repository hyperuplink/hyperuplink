package helpers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/forum"
)

type CategoryWithForums struct {
	Category category.Category
	Forums   []forum.Forum
}

func GetPaths(c fiber.Ctx) (
	pathName string,
	absPath string,
	relRoot string,
) {
	cR := c.Route()
	absPath = cR.Path
	splitAbsPath := strings.Split(absPath, "/")
	lenSplitAbsPath := len(splitAbsPath)
	if lenSplitAbsPath > 0 {
		pathName = splitAbsPath[(lenSplitAbsPath - 1)]
	} else {
		pathName = ""
	}
	parts := lenSplitAbsPath - 1

	relRoot = ""
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}

	return pathName, absPath, relRoot
}
